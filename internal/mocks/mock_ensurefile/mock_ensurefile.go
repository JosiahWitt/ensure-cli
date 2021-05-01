package mock_ensurefile

import (
	gomock "github.com/golang/mock/gomock"
	"reflect"
)

import (
	"github.com/JosiahWitt/ensure-cli/internal/ensurefile"
)

type MockLoaderIface struct {
	ctrl     *gomock.Controller
	recorder *MockLoaderIfaceMockRecorder
}

type MockLoaderIfaceMockRecorder struct {
	mock *MockLoaderIface
}

func NewMockLoaderIface(ctrl *gomock.Controller) *MockLoaderIface {
	mock := &MockLoaderIface{ctrl: ctrl}
	mock.recorder = &MockLoaderIfaceMockRecorder{mock: mock}
	return mock
}

func (m *MockLoaderIface) EXPECT() *MockLoaderIfaceMockRecorder {
	return m.recorder
}

func (m *MockLoaderIface) NEW(ctrl *gomock.Controller) *MockLoaderIface {
	return NewMockLoaderIface(ctrl)
}

func (m *MockLoaderIface) LoadConfig(_pwd string) (*ensurefile.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadConfig", _pwd)
	ret0, _ := ret[0].(*ensurefile.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockLoaderIfaceMockRecorder) LoadConfig(_pwd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadConfig", reflect.TypeOf((*MockLoaderIface)(nil).LoadConfig), _pwd)
}
