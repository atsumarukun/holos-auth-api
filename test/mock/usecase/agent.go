// Code generated by MockGen. DO NOT EDIT.
// Source: agent.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	apierr "holos-auth-api/internal/app/api/pkg/apierr"
	dto "holos-auth-api/internal/app/api/usecase/dto"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockAgentUsecase is a mock of AgentUsecase interface.
type MockAgentUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockAgentUsecaseMockRecorder
}

// MockAgentUsecaseMockRecorder is the mock recorder for MockAgentUsecase.
type MockAgentUsecaseMockRecorder struct {
	mock *MockAgentUsecase
}

// NewMockAgentUsecase creates a new mock instance.
func NewMockAgentUsecase(ctrl *gomock.Controller) *MockAgentUsecase {
	mock := &MockAgentUsecase{ctrl: ctrl}
	mock.recorder = &MockAgentUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAgentUsecase) EXPECT() *MockAgentUsecaseMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockAgentUsecase) Create(arg0 context.Context, arg1 uuid.UUID, arg2 string) (*dto.AgentDTO, apierr.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2)
	ret0, _ := ret[0].(*dto.AgentDTO)
	ret1, _ := ret[1].(apierr.ApiError)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockAgentUsecaseMockRecorder) Create(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAgentUsecase)(nil).Create), arg0, arg1, arg2)
}

// Delete mocks base method.
func (m *MockAgentUsecase) Delete(arg0 context.Context, arg1, arg2 uuid.UUID) apierr.ApiError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1, arg2)
	ret0, _ := ret[0].(apierr.ApiError)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockAgentUsecaseMockRecorder) Delete(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAgentUsecase)(nil).Delete), arg0, arg1, arg2)
}

// GetPolicies mocks base method.
func (m *MockAgentUsecase) GetPolicies(arg0 context.Context, arg1, arg2 uuid.UUID) ([]*dto.PolicyDTO, apierr.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPolicies", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*dto.PolicyDTO)
	ret1, _ := ret[1].(apierr.ApiError)
	return ret0, ret1
}

// GetPolicies indicates an expected call of GetPolicies.
func (mr *MockAgentUsecaseMockRecorder) GetPolicies(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPolicies", reflect.TypeOf((*MockAgentUsecase)(nil).GetPolicies), arg0, arg1, arg2)
}

// Gets mocks base method.
func (m *MockAgentUsecase) Gets(arg0 context.Context, arg1 uuid.UUID) ([]*dto.AgentDTO, apierr.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Gets", arg0, arg1)
	ret0, _ := ret[0].([]*dto.AgentDTO)
	ret1, _ := ret[1].(apierr.ApiError)
	return ret0, ret1
}

// Gets indicates an expected call of Gets.
func (mr *MockAgentUsecaseMockRecorder) Gets(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Gets", reflect.TypeOf((*MockAgentUsecase)(nil).Gets), arg0, arg1)
}

// Update mocks base method.
func (m *MockAgentUsecase) Update(arg0 context.Context, arg1, arg2 uuid.UUID, arg3 string) (*dto.AgentDTO, apierr.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*dto.AgentDTO)
	ret1, _ := ret[1].(apierr.ApiError)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockAgentUsecaseMockRecorder) Update(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAgentUsecase)(nil).Update), arg0, arg1, arg2, arg3)
}

// UpdatePolicies mocks base method.
func (m *MockAgentUsecase) UpdatePolicies(arg0 context.Context, arg1, arg2 uuid.UUID, arg3 []uuid.UUID) ([]*dto.PolicyDTO, apierr.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePolicies", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]*dto.PolicyDTO)
	ret1, _ := ret[1].(apierr.ApiError)
	return ret0, ret1
}

// UpdatePolicies indicates an expected call of UpdatePolicies.
func (mr *MockAgentUsecaseMockRecorder) UpdatePolicies(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePolicies", reflect.TypeOf((*MockAgentUsecase)(nil).UpdatePolicies), arg0, arg1, arg2, arg3)
}
