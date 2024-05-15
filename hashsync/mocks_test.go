// Code generated by MockGen. DO NOT EDIT.
// Source: ./interface.go
//
// Generated by this command:
//
//	mockgen -typed -package=hashsync -destination=./mocks_test.go -source=./interface.go
//

// Package hashsync is a generated GoMock package.
package hashsync

import (
	context "context"
	reflect "reflect"

	types "github.com/spacemeshos/go-spacemesh/common/types"
	p2p "github.com/spacemeshos/go-spacemesh/p2p"
	server "github.com/spacemeshos/go-spacemesh/p2p/server"
	gomock "go.uber.org/mock/gomock"
)

// Mockrequester is a mock of requester interface.
type Mockrequester struct {
	ctrl     *gomock.Controller
	recorder *MockrequesterMockRecorder
}

// MockrequesterMockRecorder is the mock recorder for Mockrequester.
type MockrequesterMockRecorder struct {
	mock *Mockrequester
}

// NewMockrequester creates a new mock instance.
func NewMockrequester(ctrl *gomock.Controller) *Mockrequester {
	mock := &Mockrequester{ctrl: ctrl}
	mock.recorder = &MockrequesterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockrequester) EXPECT() *MockrequesterMockRecorder {
	return m.recorder
}

// Run mocks base method.
func (m *Mockrequester) Run(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *MockrequesterMockRecorder) Run(arg0 any) *MockrequesterRunCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*Mockrequester)(nil).Run), arg0)
	return &MockrequesterRunCall{Call: call}
}

// MockrequesterRunCall wrap *gomock.Call
type MockrequesterRunCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockrequesterRunCall) Return(arg0 error) *MockrequesterRunCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockrequesterRunCall) Do(f func(context.Context) error) *MockrequesterRunCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockrequesterRunCall) DoAndReturn(f func(context.Context) error) *MockrequesterRunCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// StreamRequest mocks base method.
func (m *Mockrequester) StreamRequest(arg0 context.Context, arg1 p2p.Peer, arg2 []byte, arg3 server.StreamRequestCallback, arg4 ...string) error {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1, arg2, arg3}
	for _, a := range arg4 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "StreamRequest", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// StreamRequest indicates an expected call of StreamRequest.
func (mr *MockrequesterMockRecorder) StreamRequest(arg0, arg1, arg2, arg3 any, arg4 ...any) *MockrequesterStreamRequestCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1, arg2, arg3}, arg4...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StreamRequest", reflect.TypeOf((*Mockrequester)(nil).StreamRequest), varargs...)
	return &MockrequesterStreamRequestCall{Call: call}
}

// MockrequesterStreamRequestCall wrap *gomock.Call
type MockrequesterStreamRequestCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockrequesterStreamRequestCall) Return(arg0 error) *MockrequesterStreamRequestCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockrequesterStreamRequestCall) Do(f func(context.Context, p2p.Peer, []byte, server.StreamRequestCallback, ...string) error) *MockrequesterStreamRequestCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockrequesterStreamRequestCall) DoAndReturn(f func(context.Context, p2p.Peer, []byte, server.StreamRequestCallback, ...string) error) *MockrequesterStreamRequestCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MocksyncBase is a mock of syncBase interface.
type MocksyncBase struct {
	ctrl     *gomock.Controller
	recorder *MocksyncBaseMockRecorder
}

// MocksyncBaseMockRecorder is the mock recorder for MocksyncBase.
type MocksyncBaseMockRecorder struct {
	mock *MocksyncBase
}

// NewMocksyncBase creates a new mock instance.
func NewMocksyncBase(ctrl *gomock.Controller) *MocksyncBase {
	mock := &MocksyncBase{ctrl: ctrl}
	mock.recorder = &MocksyncBaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocksyncBase) EXPECT() *MocksyncBaseMockRecorder {
	return m.recorder
}

// count mocks base method.
func (m *MocksyncBase) count() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "count")
	ret0, _ := ret[0].(int)
	return ret0
}

// count indicates an expected call of count.
func (mr *MocksyncBaseMockRecorder) count() *MocksyncBasecountCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "count", reflect.TypeOf((*MocksyncBase)(nil).count))
	return &MocksyncBasecountCall{Call: call}
}

// MocksyncBasecountCall wrap *gomock.Call
type MocksyncBasecountCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MocksyncBasecountCall) Return(arg0 int) *MocksyncBasecountCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MocksyncBasecountCall) Do(f func() int) *MocksyncBasecountCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MocksyncBasecountCall) DoAndReturn(f func() int) *MocksyncBasecountCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// derive mocks base method.
func (m *MocksyncBase) derive(p p2p.Peer) syncer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "derive", p)
	ret0, _ := ret[0].(syncer)
	return ret0
}

// derive indicates an expected call of derive.
func (mr *MocksyncBaseMockRecorder) derive(p any) *MocksyncBasederiveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "derive", reflect.TypeOf((*MocksyncBase)(nil).derive), p)
	return &MocksyncBasederiveCall{Call: call}
}

// MocksyncBasederiveCall wrap *gomock.Call
type MocksyncBasederiveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MocksyncBasederiveCall) Return(arg0 syncer) *MocksyncBasederiveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MocksyncBasederiveCall) Do(f func(p2p.Peer) syncer) *MocksyncBasederiveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MocksyncBasederiveCall) DoAndReturn(f func(p2p.Peer) syncer) *MocksyncBasederiveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// probe mocks base method.
func (m *MocksyncBase) probe(ctx context.Context, p p2p.Peer) (ProbeResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "probe", ctx, p)
	ret0, _ := ret[0].(ProbeResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// probe indicates an expected call of probe.
func (mr *MocksyncBaseMockRecorder) probe(ctx, p any) *MocksyncBaseprobeCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "probe", reflect.TypeOf((*MocksyncBase)(nil).probe), ctx, p)
	return &MocksyncBaseprobeCall{Call: call}
}

// MocksyncBaseprobeCall wrap *gomock.Call
type MocksyncBaseprobeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MocksyncBaseprobeCall) Return(arg0 ProbeResult, arg1 error) *MocksyncBaseprobeCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MocksyncBaseprobeCall) Do(f func(context.Context, p2p.Peer) (ProbeResult, error)) *MocksyncBaseprobeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MocksyncBaseprobeCall) DoAndReturn(f func(context.Context, p2p.Peer) (ProbeResult, error)) *MocksyncBaseprobeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// run mocks base method.
func (m *MocksyncBase) run(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "run", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// run indicates an expected call of run.
func (mr *MocksyncBaseMockRecorder) run(ctx any) *MocksyncBaserunCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "run", reflect.TypeOf((*MocksyncBase)(nil).run), ctx)
	return &MocksyncBaserunCall{Call: call}
}

// MocksyncBaserunCall wrap *gomock.Call
type MocksyncBaserunCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MocksyncBaserunCall) Return(arg0 error) *MocksyncBaserunCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MocksyncBaserunCall) Do(f func(context.Context) error) *MocksyncBaserunCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MocksyncBaserunCall) DoAndReturn(f func(context.Context) error) *MocksyncBaserunCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Mocksyncer is a mock of syncer interface.
type Mocksyncer struct {
	ctrl     *gomock.Controller
	recorder *MocksyncerMockRecorder
}

// MocksyncerMockRecorder is the mock recorder for Mocksyncer.
type MocksyncerMockRecorder struct {
	mock *Mocksyncer
}

// NewMocksyncer creates a new mock instance.
func NewMocksyncer(ctrl *gomock.Controller) *Mocksyncer {
	mock := &Mocksyncer{ctrl: ctrl}
	mock.recorder = &MocksyncerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mocksyncer) EXPECT() *MocksyncerMockRecorder {
	return m.recorder
}

// peer mocks base method.
func (m *Mocksyncer) peer() p2p.Peer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "peer")
	ret0, _ := ret[0].(p2p.Peer)
	return ret0
}

// peer indicates an expected call of peer.
func (mr *MocksyncerMockRecorder) peer() *MocksyncerpeerCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "peer", reflect.TypeOf((*Mocksyncer)(nil).peer))
	return &MocksyncerpeerCall{Call: call}
}

// MocksyncerpeerCall wrap *gomock.Call
type MocksyncerpeerCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MocksyncerpeerCall) Return(arg0 p2p.Peer) *MocksyncerpeerCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MocksyncerpeerCall) Do(f func() p2p.Peer) *MocksyncerpeerCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MocksyncerpeerCall) DoAndReturn(f func() p2p.Peer) *MocksyncerpeerCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// sync mocks base method.
func (m *Mocksyncer) sync(ctx context.Context, x, y *types.Hash32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "sync", ctx, x, y)
	ret0, _ := ret[0].(error)
	return ret0
}

// sync indicates an expected call of sync.
func (mr *MocksyncerMockRecorder) sync(ctx, x, y any) *MocksyncersyncCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "sync", reflect.TypeOf((*Mocksyncer)(nil).sync), ctx, x, y)
	return &MocksyncersyncCall{Call: call}
}

// MocksyncersyncCall wrap *gomock.Call
type MocksyncersyncCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MocksyncersyncCall) Return(arg0 error) *MocksyncersyncCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MocksyncersyncCall) Do(f func(context.Context, *types.Hash32, *types.Hash32) error) *MocksyncersyncCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MocksyncersyncCall) DoAndReturn(f func(context.Context, *types.Hash32, *types.Hash32) error) *MocksyncersyncCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MocksyncRunner is a mock of syncRunner interface.
type MocksyncRunner struct {
	ctrl     *gomock.Controller
	recorder *MocksyncRunnerMockRecorder
}

// MocksyncRunnerMockRecorder is the mock recorder for MocksyncRunner.
type MocksyncRunnerMockRecorder struct {
	mock *MocksyncRunner
}

// NewMocksyncRunner creates a new mock instance.
func NewMocksyncRunner(ctrl *gomock.Controller) *MocksyncRunner {
	mock := &MocksyncRunner{ctrl: ctrl}
	mock.recorder = &MocksyncRunnerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocksyncRunner) EXPECT() *MocksyncRunnerMockRecorder {
	return m.recorder
}

// fullSync mocks base method.
func (m *MocksyncRunner) fullSync(ctx context.Context, syncPeers []p2p.Peer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "fullSync", ctx, syncPeers)
	ret0, _ := ret[0].(error)
	return ret0
}

// fullSync indicates an expected call of fullSync.
func (mr *MocksyncRunnerMockRecorder) fullSync(ctx, syncPeers any) *MocksyncRunnerfullSyncCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "fullSync", reflect.TypeOf((*MocksyncRunner)(nil).fullSync), ctx, syncPeers)
	return &MocksyncRunnerfullSyncCall{Call: call}
}

// MocksyncRunnerfullSyncCall wrap *gomock.Call
type MocksyncRunnerfullSyncCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MocksyncRunnerfullSyncCall) Return(arg0 error) *MocksyncRunnerfullSyncCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MocksyncRunnerfullSyncCall) Do(f func(context.Context, []p2p.Peer) error) *MocksyncRunnerfullSyncCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MocksyncRunnerfullSyncCall) DoAndReturn(f func(context.Context, []p2p.Peer) error) *MocksyncRunnerfullSyncCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// splitSync mocks base method.
func (m *MocksyncRunner) splitSync(ctx context.Context, syncPeers []p2p.Peer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "splitSync", ctx, syncPeers)
	ret0, _ := ret[0].(error)
	return ret0
}

// splitSync indicates an expected call of splitSync.
func (mr *MocksyncRunnerMockRecorder) splitSync(ctx, syncPeers any) *MocksyncRunnersplitSyncCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "splitSync", reflect.TypeOf((*MocksyncRunner)(nil).splitSync), ctx, syncPeers)
	return &MocksyncRunnersplitSyncCall{Call: call}
}

// MocksyncRunnersplitSyncCall wrap *gomock.Call
type MocksyncRunnersplitSyncCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MocksyncRunnersplitSyncCall) Return(arg0 error) *MocksyncRunnersplitSyncCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MocksyncRunnersplitSyncCall) Do(f func(context.Context, []p2p.Peer) error) *MocksyncRunnersplitSyncCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MocksyncRunnersplitSyncCall) DoAndReturn(f func(context.Context, []p2p.Peer) error) *MocksyncRunnersplitSyncCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
