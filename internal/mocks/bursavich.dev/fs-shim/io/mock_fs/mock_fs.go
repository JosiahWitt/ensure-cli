package mock_fs

import (
	gomock "github.com/golang/mock/gomock"
	"reflect"
)

import (
	"io/fs"
)

type MockReadFileFS struct {
	ctrl     *gomock.Controller
	recorder *MockReadFileFSMockRecorder
}

type MockReadFileFSMockRecorder struct {
	mock *MockReadFileFS
}

func NewMockReadFileFS(ctrl *gomock.Controller) *MockReadFileFS {
	mock := &MockReadFileFS{ctrl: ctrl}
	mock.recorder = &MockReadFileFSMockRecorder{mock: mock}
	return mock
}

func (m *MockReadFileFS) EXPECT() *MockReadFileFSMockRecorder {
	return m.recorder
}

func (m *MockReadFileFS) NEW(ctrl *gomock.Controller) *MockReadFileFS {
	return NewMockReadFileFS(ctrl)
}

func (m *MockReadFileFS) Open(_name string) (fs.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Open", _name)
	ret0, _ := ret[0].(fs.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockReadFileFSMockRecorder) Open(_name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*MockReadFileFS)(nil).Open), _name)
}

func (m *MockReadFileFS) ReadFile(_name string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadFile", _name)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockReadFileFSMockRecorder) ReadFile(_name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadFile", reflect.TypeOf((*MockReadFileFS)(nil).ReadFile), _name)
}
