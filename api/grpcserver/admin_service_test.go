package grpcserver

import (
	"context"
	"errors"
	"io"
	"sync/atomic"
	"testing"
	"time"

	pb "github.com/spacemeshos/api/release/go/spacemesh/v1"
	"github.com/stretchr/testify/require"

	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spacemeshos/go-spacemesh/sql"
	"github.com/spacemeshos/go-spacemesh/sql/accounts"
	"github.com/spacemeshos/go-spacemesh/sql/atxs"
	"github.com/spacemeshos/go-spacemesh/sql/statesql"
)

const snapshot uint32 = 15

func newAtx(tb testing.TB, db *statesql.Database) {
	atx := &types.ActivationTx{
		PublishEpoch:   types.EpochID(2),
		Sequence:       0,
		CommitmentATX:  &types.ATXID{1},
		NumUnits:       2,
		Coinbase:       types.Address{1, 2, 3},
		BaseTickHeight: 1111,
		TickCount:      12,
		VRFNonce:       types.VRFPostIndex(11),
	}
	atx.SetID(types.RandomATXID())
	atx.SmesherID = types.BytesToNodeID(types.RandomBytes(20))
	atx.SetReceived(time.Now().Local())
	require.NoError(tb, atxs.Add(db, atx))
}

func createMesh(tb testing.TB, db *statesql.Database) {
	for range 10 {
		newAtx(tb, db)
	}
	acct := &types.Account{
		Layer:           types.LayerID(0),
		Address:         types.Address{1, 1},
		NextNonce:       1,
		Balance:         1300,
		TemplateAddress: &types.Address{2},
		State:           []byte("state10"),
	}
	require.NoError(tb, accounts.Update(db, acct))
}

func TestAdminService_Checkpoint(t *testing.T) {
	db := statesql.InMemory()
	createMesh(t, db)
	svc := NewAdminService(db, t.TempDir(), nil)
	cfg, cleanup := launchServer(t, svc)
	t.Cleanup(cleanup)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn := dialGrpc(ctx, t, cfg)
	c := pb.NewAdminServiceClient(conn)

	stream, err := c.CheckpointStream(ctx, &pb.CheckpointStreamRequest{SnapshotLayer: snapshot})
	require.NoError(t, err)

	var chunks int
	read := func() {
		for {
			select {
			case <-ctx.Done():
				t.Fail()
			default:
				msg, err := stream.Recv()
				if errors.Is(err, io.EOF) {
					return
				}
				require.NoError(t, err)
				if len(msg.Data) == chunksize {
					chunks++
				}
			}
		}
	}
	read()
	require.NotZero(t, chunks)
}

func TestAdminService_CheckpointError(t *testing.T) {
	db := statesql.InMemory()
	svc := NewAdminService(db, t.TempDir(), nil)
	cfg, cleanup := launchServer(t, svc)
	t.Cleanup(cleanup)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn := dialGrpc(ctx, t, cfg)
	c := pb.NewAdminServiceClient(conn)

	stream, err := c.CheckpointStream(ctx, &pb.CheckpointStreamRequest{SnapshotLayer: snapshot})
	require.NoError(t, err)
	_, err = stream.Recv()
	require.ErrorContains(t, err, sql.ErrNotFound.Error())
}

func TestAdminService_Recovery(t *testing.T) {
	db := statesql.InMemory()
	recoveryCalled := atomic.Bool{}
	svc := NewAdminService(db, t.TempDir(), nil)
	svc.recover = func() { recoveryCalled.Store(true) }

	cfg, cleanup := launchServer(t, svc)
	t.Cleanup(cleanup)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn := dialGrpc(ctx, t, cfg)
	c := pb.NewAdminServiceClient(conn)

	_, err := c.Recover(ctx, &pb.RecoverRequest{})
	require.NoError(t, err)
	require.True(t, recoveryCalled.Load())
}
