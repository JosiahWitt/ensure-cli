package mock_runcmd

import (
	gomock "github.com/golang/mock/gomock"
	"reflect"
)

import (
	"context"
	"github.com/JosiahWitt/ensure-cli/internal/runcmd"
)

type MockRunnerIface struct {
	ctrl     *gomock.Controller
	recorder *MockRunnerIfaceMockRecorder
}

type MockRunnerIfaceMockRecorder struct {
	mock *MockRunnerIface
}

func NewMockRunnerIface(ctrl *gomock.Controller) *MockRunnerIface {
	mock := &MockRunnerIface{ctrl: ctrl}
	mock.recorder = &MockRunnerIfaceMockRecorder{mock: mock}
	return mock
}

func (m *MockRunnerIface) EXPECT() *MockRunnerIfaceMockRecorder {
	return m.recorder
}

func (m *MockRunnerIface) NEW(ctrl *gomock.Controller) *MockRunnerIface {
	return NewMockRunnerIface(ctrl)
}

func (m *MockRunnerIface) Exec(_ctx context.Context, _params *runcmd.ExecParams) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exec", _ctx, _params)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockRunnerIfaceMockRecorder) Exec(_ctx interface{}, _params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockRunnerIface)(nil).Exec), _ctx, _params)
}
