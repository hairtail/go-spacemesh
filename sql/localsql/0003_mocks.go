// Code generated by MockGen. DO NOT EDIT.
// Source: ./0003_migration_interfaces.go
//
// Generated by this command:
//
//	mockgen -typed -package=localsql -destination=./0003_mocks.go -source=./0003_migration_interfaces.go
//

// Package localsql is a generated GoMock package.
package localsql

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockPoetClient is a mock of PoetClient interface.
type MockPoetClient struct {
	ctrl     *gomock.Controller
	recorder *MockPoetClientMockRecorder
}

// MockPoetClientMockRecorder is the mock recorder for MockPoetClient.
type MockPoetClientMockRecorder struct {
	mock *MockPoetClient
}

// NewMockPoetClient creates a new mock instance.
func NewMockPoetClient(ctrl *gomock.Controller) *MockPoetClient {
	mock := &MockPoetClient{ctrl: ctrl}
	mock.recorder = &MockPoetClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPoetClient) EXPECT() *MockPoetClientMockRecorder {
	return m.recorder
}

// Address mocks base method.
func (m *MockPoetClient) Address() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Address")
	ret0, _ := ret[0].(string)
	return ret0
}

// Address indicates an expected call of Address.
func (mr *MockPoetClientMockRecorder) Address() *MockPoetClientAddressCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Address", reflect.TypeOf((*MockPoetClient)(nil).Address))
	return &MockPoetClientAddressCall{Call: call}
}

// MockPoetClientAddressCall wrap *gomock.Call
type MockPoetClientAddressCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockPoetClientAddressCall) Return(arg0 string) *MockPoetClientAddressCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockPoetClientAddressCall) Do(f func() string) *MockPoetClientAddressCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockPoetClientAddressCall) DoAndReturn(f func() string) *MockPoetClientAddressCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// PoetServiceID mocks base method.
func (m *MockPoetClient) PoetServiceID(ctx context.Context) []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PoetServiceID", ctx)
	ret0, _ := ret[0].([]byte)
	return ret0
}

// PoetServiceID indicates an expected call of PoetServiceID.
func (mr *MockPoetClientMockRecorder) PoetServiceID(ctx any) *MockPoetClientPoetServiceIDCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PoetServiceID", reflect.TypeOf((*MockPoetClient)(nil).PoetServiceID), ctx)
	return &MockPoetClientPoetServiceIDCall{Call: call}
}

// MockPoetClientPoetServiceIDCall wrap *gomock.Call
type MockPoetClientPoetServiceIDCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockPoetClientPoetServiceIDCall) Return(arg0 []byte) *MockPoetClientPoetServiceIDCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockPoetClientPoetServiceIDCall) Do(f func(context.Context) []byte) *MockPoetClientPoetServiceIDCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockPoetClientPoetServiceIDCall) DoAndReturn(f func(context.Context) []byte) *MockPoetClientPoetServiceIDCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}