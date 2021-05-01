package mock_exitcleanup

import (
	gomock "github.com/golang/mock/gomock"
	"reflect"
)

import (
	"context"
)

type MockExitCleaner struct {
	ctrl     *gomock.Controller
	recorder *MockExitCleanerMockRecorder
}

type MockExitCleanerMockRecorder struct {
	mock *MockExitCleaner
}

func NewMockExitCleaner(ctrl *gomock.Controller) *MockExitCleaner {
	mock := &MockExitCleaner{ctrl: ctrl}
	mock.recorder = &MockExitCleanerMockRecorder{mock: mock}
	return mock
}

func (m *MockExitCleaner) EXPECT() *MockExitCleanerMockRecorder {
	return m.recorder
}

func (m *MockExitCleaner) NEW(ctrl *gomock.Controller) *MockExitCleaner {
	return NewMockExitCleaner(ctrl)
}

func (m *MockExitCleaner) Register(_fn func() error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", _fn)
	var _ = ret
	return
}

func (mr *MockExitCleanerMockRecorder) Register(_fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockExitCleaner)(nil).Register), _fn)
}

func (m *MockExitCleaner) ToContext(_ctx context.Context) context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToContext", _ctx)
	ret0, _ := ret[0].(context.Context)
	return ret0
}

func (mr *MockExitCleanerMockRecorder) ToContext(_ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToContext", reflect.TypeOf((*MockExitCleaner)(nil).ToContext), _ctx)
}
