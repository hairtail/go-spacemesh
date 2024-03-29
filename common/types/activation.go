package types

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/spacemeshos/go-scale"
	"github.com/spacemeshos/post/shared"

	"github.com/spacemeshos/go-spacemesh/codec"
	"github.com/spacemeshos/go-spacemesh/common/types/primitive"
	"github.com/spacemeshos/go-spacemesh/common/types/wire"
	"github.com/spacemeshos/go-spacemesh/common/util"
	"github.com/spacemeshos/go-spacemesh/log"
)

//go:generate scalegen -types ATXMetadata,EpochActiveSet

// BytesToATXID is a helper to copy buffer into a ATXID.
func BytesToATXID(buf []byte) (id ATXID) {
	copy(id[:], buf)
	return id
}

type Validity int

const (
	Unknown Validity = iota
	Valid
	Invalid
)

// ATXID is a 32 byte hash used to identify an activation transaction.
type ATXID Hash32

const (
	// ATXIDSize in bytes.
	ATXIDSize = Hash32Length
)

// String implements stringer interface.
func (t ATXID) String() string {
	return t.ShortString()
}

// ShortString returns the first few characters of the ID, for logging purposes.
func (t ATXID) ShortString() string {
	return t.Hash32().ShortString()
}

// Hash32 returns the ATXID as a Hash32.
func (t ATXID) Hash32() Hash32 {
	return Hash32(t)
}

// Bytes returns the ATXID as a byte slice.
func (t ATXID) Bytes() []byte {
	return Hash32(t).Bytes()
}

// Field returns a log field. Implements the LoggableField interface.
func (t ATXID) Field() log.Field { return log.FieldNamed("atx_id", t.Hash32()) }

// EncodeScale implements scale codec interface.
func (t *ATXID) EncodeScale(e *scale.Encoder) (int, error) {
	return scale.EncodeByteArray(e, t[:])
}

// DecodeScale implements scale codec interface.
func (t *ATXID) DecodeScale(d *scale.Decoder) (int, error) {
	return scale.DecodeByteArray(d, t[:])
}

func (t *ATXID) MarshalText() ([]byte, error) {
	return util.Base64Encode(t[:]), nil
}

func (t *ATXID) UnmarshalText(buf []byte) error {
	return util.Base64Decode(t[:], buf)
}

// EmptyATXID is a canonical empty ATXID.
var EmptyATXID = ATXID{}

type ATXIDs []ATXID

// impl zap's ArrayMarshaler interface.
func (ids ATXIDs) MarshalLogArray(enc log.ArrayEncoder) error {
	for _, id := range ids {
		enc.AppendString(id.String())
	}
	return nil
}

// NIPostChallenge is the set of fields that's serialized, hashed and submitted to the PoET service
// to be included in the PoET membership proof.
type NIPostChallenge struct {
	PublishEpoch EpochID
	// Sequence number counts the number of ancestors of the ATX. It sequentially increases for each ATX in the chain.
	// Two ATXs with the same sequence number from the same miner can be used as the proof of malfeasance against
	// that miner.
	Sequence uint64
	// the previous ATX's ID (for all but the first in the sequence)
	PrevATXID      ATXID
	PositioningATX ATXID

	// CommitmentATX is the ATX used in the commitment for initializing the PoST of the node.
	CommitmentATX *ATXID
	InitialPost   *Post
}

func (c *NIPostChallenge) MarshalLogObject(encoder log.ObjectEncoder) error {
	if c == nil {
		return nil
	}
	encoder.AddUint32("PublishEpoch", c.PublishEpoch.Uint32())
	encoder.AddUint64("Sequence", c.Sequence)
	encoder.AddString("PrevATXID", c.PrevATXID.String())
	encoder.AddUint32("PublishEpoch", c.PublishEpoch.Uint32())
	encoder.AddString("PositioningATX", c.PositioningATX.String())
	if c.CommitmentATX != nil {
		encoder.AddString("CommitmentATX", c.CommitmentATX.String())
	}
	encoder.AddObject("InitialPost", c.InitialPost)
	return nil
}

// TargetEpoch returns the target epoch of the NIPostChallenge. This is the epoch in which the miner is eligible
// to participate thanks to the ATX.
func (challenge *NIPostChallenge) TargetEpoch() EpochID {
	return challenge.PublishEpoch + 1
}

func (challenge *NIPostChallenge) Hash() primitive.Hash32 {
	return challenge.ToWireV1().Hash()
}

// ATXMetadata is the data of ActivationTx that is signed.
// It is also used for Malfeasance proofs.
type ATXMetadata struct {
	PublishEpoch EpochID
	MsgHash      Hash32 // Hash of InnerActivationTx (returned by HashInnerBytes)
}

func (m *ATXMetadata) MarshalLogObject(encoder log.ObjectEncoder) error {
	encoder.AddUint32("epoch", uint32(m.PublishEpoch))
	encoder.AddString("hash", m.MsgHash.ShortString())
	return nil
}

// ActivationTx is a full, signed activation transaction. It includes (or references) everything a miner needs to prove
// they are eligible to actively participate in the Spacemesh protocol in the next epoch.
type ActivationTx struct {
	NIPostChallenge
	Coinbase Address
	NumUnits uint32

	NIPost   NIPost
	VRFNonce *VRFPostIndex

	SmesherID NodeID
	Signature EdSignature

	golden            bool
	id                ATXID     // non-exported cache of the ATXID
	effectiveNumUnits uint32    // the number of effective units in the ATX (minimum of this ATX and the previous ATX)
	received          time.Time // time received by node, gossiped or synced
	validity          Validity  // whether the chain is fully verified and OK
}

// NewActivationTx returns a new activation transaction. The ATXID is calculated and cached.
func NewActivationTx(
	challenge NIPostChallenge,
	coinbase Address,
	nipost NIPost,
	numUnits uint32,
	nonce *VRFPostIndex,
) *ActivationTx {
	atx := &ActivationTx{
		NIPostChallenge: challenge,
		Coinbase:        coinbase,
		NumUnits:        numUnits,
		NIPost:          nipost,
		VRFNonce:        nonce,
	}
	return atx
}

// Golden returns true if atx is from a checkpoint snapshot.
// a golden ATX is not verifiable, and is only allowed to be prev atx or positioning atx.
func (atx *ActivationTx) Golden() bool {
	return atx.golden
}

// SetGolden set atx to golden.
func (atx *ActivationTx) SetGolden() {
	atx.golden = true
}

// MarshalLogObject implements logging interface.
func (atx *ActivationTx) MarshalLogObject(encoder log.ObjectEncoder) error {
	encoder.AddString("atx_id", atx.id.String())
	encoder.AddString("challenge", atx.NIPostChallenge.Hash().String())
	encoder.AddString("smesher", atx.SmesherID.String())
	encoder.AddString("prev_atx_id", atx.PrevATXID.String())
	encoder.AddString("pos_atx_id", atx.PositioningATX.String())
	if atx.CommitmentATX != nil {
		encoder.AddString("commitment_atx_id", atx.CommitmentATX.String())
	}
	if atx.VRFNonce != nil {
		encoder.AddUint64("vrf_nonce", uint64(*atx.VRFNonce))
	}
	encoder.AddString("coinbase", atx.Coinbase.String())
	encoder.AddUint32("epoch", atx.PublishEpoch.Uint32())
	encoder.AddUint64("num_units", uint64(atx.NumUnits))
	if atx.effectiveNumUnits != 0 {
		encoder.AddUint64("effective_num_units", uint64(atx.effectiveNumUnits))
	}
	encoder.AddUint64("sequence_number", atx.Sequence)
	return nil
}

// Initialize calculates and sets the cached ID field. This field must be set before calling the ID() method.
func (atx *ActivationTx) Initialize() error {
	if atx.ID() != EmptyATXID {
		return errors.New("ATX already initialized")
	}

	atx.id = ATXID(Hash32(atx.ToWireV1().HashInnerBytes()))
	return nil
}

// GetPoetProofRef returns the reference to the PoET proof.
func (atx *ActivationTx) GetPoetProofRef() Hash32 {
	return BytesToHash(atx.NIPost.PostMetadata.Challenge)
}

// ShortString returns the first 5 characters of the ID, for logging purposes.
func (atx *ActivationTx) ShortString() string {
	return atx.ID().ShortString()
}

// ID returns the ATX's ID.
func (atx *ActivationTx) ID() ATXID {
	return atx.id
}

func (atx *ActivationTx) EffectiveNumUnits() uint32 {
	if atx.effectiveNumUnits == 0 {
		panic("effectiveNumUnits field must be set")
	}
	return atx.effectiveNumUnits
}

// SetID sets the ATXID in this ATX's cache.
func (atx *ActivationTx) SetID(id ATXID) {
	atx.id = id
}

func (atx *ActivationTx) SetEffectiveNumUnits(numUnits uint32) {
	atx.effectiveNumUnits = numUnits
}

func (atx *ActivationTx) SetReceived(received time.Time) {
	atx.received = received
}

func (atx *ActivationTx) Received() time.Time {
	return atx.received
}

func (atx *ActivationTx) Validity() Validity {
	return atx.validity
}

func (atx *ActivationTx) SetValidity(validity Validity) {
	atx.validity = validity
}

// Verify an ATX for a given base TickHeight and TickCount.
func (atx *ActivationTx) Verify(baseTickHeight, tickCount uint64) (*VerifiedActivationTx, error) {
	if atx.id == EmptyATXID {
		if err := atx.Initialize(); err != nil {
			return nil, err
		}
	}
	if atx.effectiveNumUnits == 0 {
		return nil, errors.New("effective num units not set")
	}
	if !atx.Golden() && atx.received.IsZero() {
		return nil, errors.New("received time not set")
	}
	vAtx := &VerifiedActivationTx{
		ActivationTx: atx,

		baseTickHeight: baseTickHeight,
		tickCount:      tickCount,
	}
	return vAtx, nil
}

// Merkle proof proving that a given leaf is included in the root of merkle tree.
type MerkleProof struct {
	// Nodes on path from leaf to root (not including leaf)
	Nodes     []Hash32 `scale:"max=32"`
	LeafIndex uint64
}

// NIPost is Non-Interactive Proof of Space-Time.
// Given an id, a space parameter S, a duration D and a challenge C,
// it can convince a verifier that (1) the prover expended S * D space-time
// after learning the challenge C. (2) the prover did not know the NIPost until D time
// after the prover learned C.
type NIPost struct {
	// Membership proves that the challenge for the PoET, which is
	// constructed from fields in the activation transaction,
	// is a member of the poet's proof.
	// Proof.Root must match the Poet's POSW statement.
	Membership MerkleProof

	// Post is the proof that the prover data is still stored (or was recomputed) at
	// the time he learned the challenge constructed from the PoET.
	Post Post

	// PostMetadata is the Post metadata, associated with the proof.
	// The proof should be verified upon the metadata during the syntactic validation,
	// while the metadata should be verified during the contextual validation.
	PostMetadata PostMetadata
}

// VRFPostIndex is the nonce generated using Pow during post initialization. It is used as a mitigation for
// grinding of identities for VRF eligibility.
type VRFPostIndex uint64

// Field returns a log field. Implements the LoggableField interface.
func (v VRFPostIndex) Field() log.Field { return log.Uint64("vrf_nonce", uint64(v)) }

// Post is an alias to postShared.Proof.
type Post shared.Proof

func (p *Post) MarshalLogObject(encoder log.ObjectEncoder) error {
	if p == nil {
		return nil
	}
	encoder.AddUint32("nonce", p.Nonce)
	encoder.AddString("indices", hex.EncodeToString(p.Indices))
	return nil
}

// PostMetadata is similar postShared.ProofMetadata, but without the fields which can be derived elsewhere
// in a given ATX (eg. NodeID, NumUnits).
type PostMetadata struct {
	Challenge     []byte `scale:"max=32"`
	LabelsPerUnit uint64
}

func (m *PostMetadata) MarshalLogObject(encoder log.ObjectEncoder) error {
	if m == nil {
		return nil
	}
	encoder.AddString("Challenge", hex.EncodeToString(m.Challenge))
	encoder.AddUint64("LabelsPerUnit", m.LabelsPerUnit)
	return nil
}

// ToATXIDs returns a slice of ATXID corresponding to the given activation tx.
func ToATXIDs(atxs []*ActivationTx) []ATXID {
	ids := make([]ATXID, 0, len(atxs))
	for _, atx := range atxs {
		ids = append(ids, atx.ID())
	}
	return ids
}

// ATXIDsToHashes turns a list of ATXID into their Hash32 representation.
func ATXIDsToHashes(ids []ATXID) []Hash32 {
	hashes := make([]Hash32, 0, len(ids))
	for _, id := range ids {
		hashes = append(hashes, id.Hash32())
	}
	return hashes
}

type EpochActiveSet struct {
	Epoch EpochID
	Set   []ATXID `scale:"max=2700000"` // to be in line with `EpochData` in fetch/wire_types.go
}

var MaxEpochActiveSetSize = scale.MustGetMaxElements[EpochActiveSet]("Set")

func (p *Post) ToWireV1() *wire.PostV1 {
	if p == nil {
		return nil
	}
	return &wire.PostV1{
		Nonce:   p.Nonce,
		Indices: p.Indices,
		Pow:     p.Pow,
	}
}

func (n *NIPost) ToWireV1() *wire.NIPostV1 {
	if n == nil {
		return nil
	}

	return &wire.NIPostV1{
		Membership: *n.Membership.ToWireV1(),
		Post:       n.Post.ToWireV1(),
		PostMetadata: &wire.PostMetadataV1{
			Challenge:     n.PostMetadata.Challenge,
			LabelsPerUnit: n.PostMetadata.LabelsPerUnit,
		},
	}
}

func (p *MerkleProof) ToWireV1() *wire.MerkleProofV1 {
	if p == nil {
		return nil
	}
	nodes := make([]Hash32, 0, len(p.Nodes))
	for _, node := range p.Nodes {
		nodes = append(nodes, Hash32(node))
	}
	return &wire.MerkleProofV1{
		Nodes:     nodes,
		LeafIndex: p.LeafIndex,
	}
}

func (c *NIPostChallenge) ToWireV1() *wire.NIPostChallengeV1 {
	return &wire.NIPostChallengeV1{
		PublishEpoch:   c.PublishEpoch.Uint32(),
		Sequence:       c.Sequence,
		PrevATXID:      primitive.Hash32(c.PrevATXID),
		PositioningATX: primitive.Hash32(c.PositioningATX),
		CommitmentATX:  (*primitive.Hash32)(c.CommitmentATX),
		InitialPost:    c.InitialPost.ToWireV1(),
	}
}

func (a *ActivationTx) ToWireV1() *wire.ActivationTxV1 {
	atxV1 := wire.ActivationTxV1{
		InnerActivationTxV1: wire.InnerActivationTxV1{
			NIPostChallengeV1: *a.NIPostChallenge.ToWireV1(),
			Coinbase:          wire.Address(a.Coinbase),
			NumUnits:          a.NumUnits,
			NIPost:            a.NIPost.ToWireV1(),
			VRFNonce:          (*wire.VRFPostIndex)(a.VRFNonce),
		},
		SmesherID: Hash32(a.SmesherID),
		Signature: a.Signature,
	}
	if a.PrevATXID == EmptyATXID {
		atxV1.InnerActivationTxV1.NodeID = &atxV1.SmesherID
	}
	return &atxV1
}

// Decode ActivationTx from bytes.
// In future it should decide which version of ActivationTx to decode based on the publish epoch.
func AcivationTxFromBytes(data []byte) (*ActivationTx, error) {
	var wireAtx wire.ActivationTxV1
	err := codec.Decode(data, &wireAtx)
	if err != nil {
		return nil, fmt.Errorf("decoding ATX: %w", err)
	}

	return ActivationTxFromWireV1(&wireAtx)
}

func ActivationTxFromWireV1(atx *wire.ActivationTxV1) (*ActivationTx, error) {
	if (atx.PrevATXID == Hash32{}) {
		if atx.InnerActivationTxV1.NodeID == nil {
			return nil, errors.New("nil NodeID in initial ATX")
		}
	} else {
		if atx.InnerActivationTxV1.NodeID != nil {
			return nil, errors.New("non-nil NodeID in non-initial ATX")
		}
	}
	nipost, err := NIPostFromWireV1(atx.NIPost)
	if err != nil {
		return nil, err
	}

	return &ActivationTx{
		NIPostChallenge: NIPostChallenge{
			PublishEpoch:   EpochID(atx.PublishEpoch),
			Sequence:       atx.Sequence,
			PrevATXID:      ATXID(atx.PrevATXID),
			PositioningATX: ATXID(atx.PositioningATX),
			CommitmentATX:  (*ATXID)(atx.CommitmentATX),
			InitialPost:    PostFromWireV1(atx.InitialPost),
		},
		Coinbase:  Address(atx.Coinbase),
		NumUnits:  atx.NumUnits,
		NIPost:    *nipost,
		VRFNonce:  (*VRFPostIndex)(atx.VRFNonce),
		SmesherID: NodeID(atx.SmesherID),
		Signature: atx.Signature,
	}, nil
}

func NIPostChallengeFromWireV1(ch wire.NIPostChallengeV1) *NIPostChallenge {
	return &NIPostChallenge{
		PublishEpoch:   EpochID(ch.PublishEpoch),
		Sequence:       ch.Sequence,
		PrevATXID:      ATXID(ch.PrevATXID),
		PositioningATX: ATXID(ch.PositioningATX),
		CommitmentATX:  (*ATXID)(ch.CommitmentATX),
		InitialPost:    PostFromWireV1(ch.InitialPost),
	}
}

func NIPostFromWireV1(nipost *wire.NIPostV1) (*NIPost, error) {
	if nipost == nil {
		return nil, errors.New("nil nipost")
	}
	if nipost.Post == nil {
		return nil, errors.New("nil nipost.post")
	}
	if nipost.PostMetadata == nil {
		return nil, errors.New("nil nipost.postmetadata")
	}
	return &NIPost{
		Membership: *MerkleProofFromWireV1(nipost.Membership),
		Post:       *PostFromWireV1(nipost.Post),
		PostMetadata: PostMetadata{
			Challenge:     nipost.PostMetadata.Challenge,
			LabelsPerUnit: nipost.PostMetadata.LabelsPerUnit,
		},
	}, nil
}

func PostFromWireV1(post *wire.PostV1) *Post {
	if post == nil {
		return nil
	}
	return &Post{
		Nonce:   post.Nonce,
		Indices: post.Indices,
		Pow:     post.Pow,
	}
}

func MerkleProofFromWireV1(proofV1 wire.MerkleProofV1) *MerkleProof {
	proof := &MerkleProof{
		LeafIndex: proofV1.LeafIndex,
	}
	for _, node := range proofV1.Nodes {
		if proof.Nodes == nil {
			proof.Nodes = make([]Hash32, 0, len(proofV1.Nodes))
		}
		proof.Nodes = append(proof.Nodes, Hash32(node))
	}
	return proof
}
