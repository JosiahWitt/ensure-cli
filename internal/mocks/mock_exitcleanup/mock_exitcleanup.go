// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/JosiahWitt/ensure-cli/internal/exitcleanup (interfaces: ExitCleaner)

// Package mock_exitcleanup is a generated GoMock package.
package mock_exitcleanup

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockExitCleaner is a mock of ExitCleaner interface.
type MockExitCleaner struct {
	ctrl     *gomock.Controller
	recorder *MockExitCleanerMockRecorder
}

// MockExitCleanerMockRecorder is the mock recorder for MockExitCleaner.
type MockExitCleanerMockRecorder struct {
	mock *MockExitCleaner
}

// NewMockExitCleaner creates a new mock instance.
func NewMockExitCleaner(ctrl *gomock.Controller) *MockExitCleaner {
	mock := &MockExitCleaner{ctrl: ctrl}
	mock.recorder = &MockExitCleanerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExitCleaner) EXPECT() *MockExitCleanerMockRecorder {
	return m.recorder
}

// Register mocks base method.
func (m *MockExitCleaner) Register(arg0 func() error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Register", arg0)
}

// Register indicates an expected call of Register.
func (mr *MockExitCleanerMockRecorder) Register(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockExitCleaner)(nil).Register), arg0)
}

// ToContext mocks base method.
func (m *MockExitCleaner) ToContext(arg0 context.Context) context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToContext", arg0)
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// ToContext indicates an expected call of ToContext.
func (mr *MockExitCleanerMockRecorder) ToContext(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToContext", reflect.TypeOf((*MockExitCleaner)(nil).ToContext), arg0)
}

// NEW creates a MockExitCleaner.
func (*MockExitCleaner) NEW(ctrl *gomock.Controller) *MockExitCleaner {
	return NewMockExitCleaner(ctrl)
}
