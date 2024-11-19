// Code generated by MockGen. DO NOT EDIT.
// Source: agent.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	entity "holos-auth-api/internal/app/api/domain/entity"
	apierr "holos-auth-api/internal/app/api/pkg/apierr"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAgentService is a mock of AgentService interface.
type MockAgentService struct {
	ctrl     *gomock.Controller
	recorder *MockAgentServiceMockRecorder
}

// MockAgentServiceMockRecorder is the mock recorder for MockAgentService.
type MockAgentServiceMockRecorder struct {
	mock *MockAgentService
}

// NewMockAgentService creates a new mock instance.
func NewMockAgentService(ctrl *gomock.Controller) *MockAgentService {
	mock := &MockAgentService{ctrl: ctrl}
	mock.recorder = &MockAgentServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAgentService) EXPECT() *MockAgentServiceMockRecorder {
	return m.recorder
}

// Exists mocks base method.
func (m *MockAgentService) Exists(arg0 context.Context, arg1 *entity.Agent) (bool, apierr.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(apierr.ApiError)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *MockAgentServiceMockRecorder) Exists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockAgentService)(nil).Exists), arg0, arg1)
}
