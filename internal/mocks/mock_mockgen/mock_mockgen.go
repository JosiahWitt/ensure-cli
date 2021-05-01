package mock_mockgen

import (
	gomock "github.com/golang/mock/gomock"
	"reflect"
)

import (
	"context"
	"github.com/JosiahWitt/ensure-cli/internal/ensurefile"
)

type MockMockGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockMockGeneratorMockRecorder
}

type MockMockGeneratorMockRecorder struct {
	mock *MockMockGenerator
}

func NewMockMockGenerator(ctrl *gomock.Controller) *MockMockGenerator {
	mock := &MockMockGenerator{ctrl: ctrl}
	mock.recorder = &MockMockGeneratorMockRecorder{mock: mock}
	return mock
}

func (m *MockMockGenerator) EXPECT() *MockMockGeneratorMockRecorder {
	return m.recorder
}

func (m *MockMockGenerator) NEW(ctrl *gomock.Controller) *MockMockGenerator {
	return NewMockMockGenerator(ctrl)
}

func (m *MockMockGenerator) GenerateMocks(_ctx context.Context, _config *ensurefile.Config) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateMocks", _ctx, _config)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockMockGeneratorMockRecorder) GenerateMocks(_ctx interface{}, _config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateMocks", reflect.TypeOf((*MockMockGenerator)(nil).GenerateMocks), _ctx, _config)
}

func (m *MockMockGenerator) TidyMocks(_config *ensurefile.Config) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TidyMocks", _config)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockMockGeneratorMockRecorder) TidyMocks(_config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TidyMocks", reflect.TypeOf((*MockMockGenerator)(nil).TidyMocks), _config)
}
