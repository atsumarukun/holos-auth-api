// Code generated by MockGen. DO NOT EDIT.
// Source: policy.go

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

// MockPolicyUsecase is a mock of PolicyUsecase interface.
type MockPolicyUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockPolicyUsecaseMockRecorder
}

// MockPolicyUsecaseMockRecorder is the mock recorder for MockPolicyUsecase.
type MockPolicyUsecaseMockRecorder struct {
	mock *MockPolicyUsecase
}

// NewMockPolicyUsecase creates a new mock instance.
func NewMockPolicyUsecase(ctrl *gomock.Controller) *MockPolicyUsecase {
	mock := &MockPolicyUsecase{ctrl: ctrl}
	mock.recorder = &MockPolicyUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPolicyUsecase) EXPECT() *MockPolicyUsecaseMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPolicyUsecase) Create(arg0 context.Context, arg1 uuid.UUID, arg2, arg3, arg4 string, arg5 []string) (*dto.PolicyDTO, apierr.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(*dto.PolicyDTO)
	ret1, _ := ret[1].(apierr.ApiError)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockPolicyUsecaseMockRecorder) Create(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPolicyUsecase)(nil).Create), arg0, arg1, arg2, arg3, arg4, arg5)
}

// Delete mocks base method.
func (m *MockPolicyUsecase) Delete(arg0 context.Context, arg1, arg2 uuid.UUID) apierr.ApiError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1, arg2)
	ret0, _ := ret[0].(apierr.ApiError)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockPolicyUsecaseMockRecorder) Delete(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockPolicyUsecase)(nil).Delete), arg0, arg1, arg2)
}

// GetAgents mocks base method.
func (m *MockPolicyUsecase) GetAgents(arg0 context.Context, arg1, arg2 uuid.UUID) ([]*dto.AgentDTO, apierr.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAgents", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*dto.AgentDTO)
	ret1, _ := ret[1].(apierr.ApiError)
	return ret0, ret1
}

// GetAgents indicates an expected call of GetAgents.
func (mr *MockPolicyUsecaseMockRecorder) GetAgents(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAgents", reflect.TypeOf((*MockPolicyUsecase)(nil).GetAgents), arg0, arg1, arg2)
}

// Gets mocks base method.
func (m *MockPolicyUsecase) Gets(arg0 context.Context, arg1 uuid.UUID) ([]*dto.PolicyDTO, apierr.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Gets", arg0, arg1)
	ret0, _ := ret[0].([]*dto.PolicyDTO)
	ret1, _ := ret[1].(apierr.ApiError)
	return ret0, ret1
}

// Gets indicates an expected call of Gets.
func (mr *MockPolicyUsecaseMockRecorder) Gets(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Gets", reflect.TypeOf((*MockPolicyUsecase)(nil).Gets), arg0, arg1)
}

// Update mocks base method.
func (m *MockPolicyUsecase) Update(arg0 context.Context, arg1, arg2 uuid.UUID, arg3, arg4, arg5 string, arg6 []string) (*dto.PolicyDTO, apierr.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2, arg3, arg4, arg5, arg6)
	ret0, _ := ret[0].(*dto.PolicyDTO)
	ret1, _ := ret[1].(apierr.ApiError)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockPolicyUsecaseMockRecorder) Update(arg0, arg1, arg2, arg3, arg4, arg5, arg6 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPolicyUsecase)(nil).Update), arg0, arg1, arg2, arg3, arg4, arg5, arg6)
}

// UpdateAgents mocks base method.
func (m *MockPolicyUsecase) UpdateAgents(arg0 context.Context, arg1, arg2 uuid.UUID, arg3 []uuid.UUID) ([]*dto.AgentDTO, apierr.ApiError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAgents", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]*dto.AgentDTO)
	ret1, _ := ret[1].(apierr.ApiError)
	return ret0, ret1
}

// UpdateAgents indicates an expected call of UpdateAgents.
func (mr *MockPolicyUsecaseMockRecorder) UpdateAgents(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAgents", reflect.TypeOf((*MockPolicyUsecase)(nil).UpdateAgents), arg0, arg1, arg2, arg3)
}
