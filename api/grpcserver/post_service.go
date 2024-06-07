package grpcserver

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/spacemeshos/api/release/go/spacemesh/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/spacemeshos/go-spacemesh/activation"
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spacemeshos/go-spacemesh/signing"
)

// PostService is a grpc server that PoST nodes can connect to in order to register.
// The bidirectional stream established between the node and the PoST node can be used
// to send challenges and receive proofs.
type PostService struct {
	root_dir string
	log      *zap.Logger

	clientMtx        sync.Mutex
	allowConnections bool
	client           map[types.NodeID]*postClient
}

type postCommand struct {
	req  *pb.NodeRequest
	resp chan<- *pb.ServiceResponse
}

func SaveNodeKey(path string, key []byte) error {
	key_path := filepath.Join(path, string(key[:]))
	err := os.WriteFile(key_path, key, 0o600)
	if err != nil {
		fmt.Printf("failed to write identity file: %v \n", key_path)
		return fmt.Errorf("write to disk failure: %w", err)
	}
	return nil
}

func LoadNodeKey(path string, key []byte, GenesisID []byte) (*signing.EdSigner, error) {
	key_path := filepath.Join(path, string(key[:]))

	data, err := os.ReadFile(key_path)
	if err != nil {
		return nil, fmt.Errorf("failed to read node key")
	}
	dst := make([]byte, signing.PrivateKeySize)
	n, err := hex.Decode(dst, data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode node key")
	}
	if n != signing.PrivateKeySize {
		return nil, fmt.Errorf("failed to parse node key size")
	}
	edSgn, err := signing.NewEdSigner(
		signing.WithPrivateKey(dst),
		signing.WithPrefix(GenesisID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to construct signer")
	}
	return edSgn, nil
}

// RegisterService registers this service with a grpc server instance.
func (s *PostService) RegisterService(server *grpc.Server) {
	pb.RegisterPostServiceServer(server, s)
}

func (s *PostService) RegisterHandlerService(mux *runtime.ServeMux) error {
	return pb.RegisterPostServiceHandlerServer(context.Background(), mux, s)
}

// String returns the name of this service.
func (s *PostService) String() string {
	return "PostService"
}

// NewPostService creates a new instance of the post grpc service.
func NewPostService(log *zap.Logger, root_dir string) *PostService {
	return &PostService{
		root_dir: root_dir,
		log:      log,
		client:   make(map[types.NodeID]*postClient),
	}
}

// AllowConnections sets if the grpc service accepts new incoming connections from post services.
func (s *PostService) AllowConnections(allow bool) {
	s.clientMtx.Lock()
	defer s.clientMtx.Unlock()
	s.allowConnections = allow
}

// connectionAllowed returns if the grpc service accepts new incoming connections from post services.
func (s *PostService) connectionAllowed() bool {
	s.clientMtx.Lock()
	defer s.clientMtx.Unlock()
	return s.allowConnections
}

// Register is called by the PoST service to connect with the node.
// It creates a bidirectional stream that is kept open until either side closes it.
// The other functions on this service are called by services of the node to send
// requests to the PoST node and receive responses.
func (s *PostService) Register(stream pb.PostService_RegisterServer) error {
	if !s.connectionAllowed() {
		return status.Error(codes.FailedPrecondition, "connection not allowed: node has no coinbase set in config")
	}

	err := stream.SendMsg(&pb.NodeRequest{
		Kind: &pb.NodeRequest_Metadata{
			Metadata: &pb.MetadataRequest{},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send metadata request: %w", err)
	}
	resp, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("failed to receive response: %w", err)
	}
	metadataResp := resp.GetMetadata()
	if metadataResp == nil {
		return fmt.Errorf("expected metadata response, got %T", resp.Kind)
	}
	meta := metadataResp.GetMeta()
	if meta == nil {
		return errors.New("expected metadata, got empty response")
	}

	// todo: add Key to meta struct; save Key in post server; register to atxBuilder.

	// if err := SaveNodeKey(s.root_dir, meta.Key); err != nil {
	// 	fmt.Printf("failed to save node key, %v \n", string(meta.Key[:]))
	// }

	// if signer, err := LoadNodeKey(s.root_dir, meta.Key, s.genesis.Bytes()); err != nil {
	// 	fmt.Printf("failed to save node key, %v \n", string(meta.Key[:]))
	// } else {
	// 	s.atxBuilder.Register(signer)
	// }

	con := make(chan postCommand)
	if err := s.setConnection(types.BytesToNodeID(meta.NodeId), con); err != nil {
		return err
	}
	defer s.dropConnection(types.BytesToNodeID(meta.NodeId))

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case cmd := <-con:
			if err := stream.SendMsg(cmd.req); err != nil {
				s.log.Error("failed to send request", zap.Error(err))
				continue
			}

			resp, err := stream.Recv()
			if err != nil {
				return fmt.Errorf("failed to receive response: %w", err)
			}

			cmd.resp <- resp
		}
	}
}

func (s *PostService) setConnection(nodeId types.NodeID, con chan postCommand) error {
	s.clientMtx.Lock()
	defer s.clientMtx.Unlock()

	if _, ok := s.client[nodeId]; ok {
		return errors.New("post service already registered")
	}
	s.client[nodeId] = newPostClient(con)
	s.log.Info("post service registered", zap.Stringer("node_id", nodeId))
	return nil
}

func (s *PostService) dropConnection(nodeId types.NodeID) error {
	s.clientMtx.Lock()
	defer s.clientMtx.Unlock()

	err := s.client[nodeId].Close()
	delete(s.client, nodeId)
	return err
}

func (s *PostService) Client(nodeId types.NodeID) (activation.PostClient, error) {
	s.clientMtx.Lock()
	defer s.clientMtx.Unlock()

	client, ok := s.client[nodeId]
	if !ok {
		return nil, activation.ErrPostClientNotConnected
	}

	return client, nil
}
