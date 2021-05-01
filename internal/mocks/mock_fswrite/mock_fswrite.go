package mock_fswrite

import (
	gomock "github.com/golang/mock/gomock"
	"reflect"
)

import (
	"io/fs"
)

type MockFSWriteIface struct {
	ctrl     *gomock.Controller
	recorder *MockFSWriteIfaceMockRecorder
}

type MockFSWriteIfaceMockRecorder struct {
	mock *MockFSWriteIface
}

func NewMockFSWriteIface(ctrl *gomock.Controller) *MockFSWriteIface {
	mock := &MockFSWriteIface{ctrl: ctrl}
	mock.recorder = &MockFSWriteIfaceMockRecorder{mock: mock}
	return mock
}

func (m *MockFSWriteIface) EXPECT() *MockFSWriteIfaceMockRecorder {
	return m.recorder
}

func (m *MockFSWriteIface) NEW(ctrl *gomock.Controller) *MockFSWriteIface {
	return NewMockFSWriteIface(ctrl)
}

func (m *MockFSWriteIface) GlobRemoveAll(_pattern string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GlobRemoveAll", _pattern)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockFSWriteIfaceMockRecorder) GlobRemoveAll(_pattern interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GlobRemoveAll", reflect.TypeOf((*MockFSWriteIface)(nil).GlobRemoveAll), _pattern)
}

func (m *MockFSWriteIface) ListRecursive(_dir string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRecursive", _dir)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockFSWriteIfaceMockRecorder) ListRecursive(_dir interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRecursive", reflect.TypeOf((*MockFSWriteIface)(nil).ListRecursive), _dir)
}

func (m *MockFSWriteIface) MkdirAll(_path string, _perm fs.FileMode) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MkdirAll", _path, _perm)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockFSWriteIfaceMockRecorder) MkdirAll(_path interface{}, _perm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MkdirAll", reflect.TypeOf((*MockFSWriteIface)(nil).MkdirAll), _path, _perm)
}

func (m *MockFSWriteIface) RemoveAll(_paths string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveAll", _paths)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockFSWriteIfaceMockRecorder) RemoveAll(_paths interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAll", reflect.TypeOf((*MockFSWriteIface)(nil).RemoveAll), _paths)
}

func (m *MockFSWriteIface) WriteFile(_filename string, _data string, _perm fs.FileMode) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteFile", _filename, _data, _perm)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockFSWriteIfaceMockRecorder) WriteFile(_filename interface{}, _data interface{}, _perm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteFile", reflect.TypeOf((*MockFSWriteIface)(nil).WriteFile), _filename, _data, _perm)
}
