// Code generated by MockGen. DO NOT EDIT.
// Source: agent.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	entity "holos-auth-api/internal/app/api/domain/entity"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockAgentRepository is a mock of AgentRepository interface.
type MockAgentRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAgentRepositoryMockRecorder
}

// MockAgentRepositoryMockRecorder is the mock recorder for MockAgentRepository.
type MockAgentRepositoryMockRecorder struct {
	mock *MockAgentRepository
}

// NewMockAgentRepository creates a new mock instance.
func NewMockAgentRepository(ctrl *gomock.Controller) *MockAgentRepository {
	mock := &MockAgentRepository{ctrl: ctrl}
	mock.recorder = &MockAgentRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAgentRepository) EXPECT() *MockAgentRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockAgentRepository) Create(arg0 context.Context, arg1 *entity.Agent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockAgentRepositoryMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAgentRepository)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockAgentRepository) Delete(arg0 context.Context, arg1 *entity.Agent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockAgentRepositoryMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAgentRepository)(nil).Delete), arg0, arg1)
}

// FindByIDsAndUserIDAndNotDeleted mocks base method.
func (m *MockAgentRepository) FindByIDsAndUserIDAndNotDeleted(arg0 context.Context, arg1 []uuid.UUID, arg2 uuid.UUID) ([]*entity.Agent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByIDsAndUserIDAndNotDeleted", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*entity.Agent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByIDsAndUserIDAndNotDeleted indicates an expected call of FindByIDsAndUserIDAndNotDeleted.
func (mr *MockAgentRepositoryMockRecorder) FindByIDsAndUserIDAndNotDeleted(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByIDsAndUserIDAndNotDeleted", reflect.TypeOf((*MockAgentRepository)(nil).FindByIDsAndUserIDAndNotDeleted), arg0, arg1, arg2)
}

// FindByNamePrefixAndUserIDAndNotDeleted mocks base method.
func (m *MockAgentRepository) FindByNamePrefixAndUserIDAndNotDeleted(arg0 context.Context, arg1 string, arg2 uuid.UUID) ([]*entity.Agent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByNamePrefixAndUserIDAndNotDeleted", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*entity.Agent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByNamePrefixAndUserIDAndNotDeleted indicates an expected call of FindByNamePrefixAndUserIDAndNotDeleted.
func (mr *MockAgentRepositoryMockRecorder) FindByNamePrefixAndUserIDAndNotDeleted(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByNamePrefixAndUserIDAndNotDeleted", reflect.TypeOf((*MockAgentRepository)(nil).FindByNamePrefixAndUserIDAndNotDeleted), arg0, arg1, arg2)
}

// FindOneByIDAndUserIDAndNotDeleted mocks base method.
func (m *MockAgentRepository) FindOneByIDAndUserIDAndNotDeleted(arg0 context.Context, arg1, arg2 uuid.UUID) (*entity.Agent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneByIDAndUserIDAndNotDeleted", arg0, arg1, arg2)
	ret0, _ := ret[0].(*entity.Agent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneByIDAndUserIDAndNotDeleted indicates an expected call of FindOneByIDAndUserIDAndNotDeleted.
func (mr *MockAgentRepositoryMockRecorder) FindOneByIDAndUserIDAndNotDeleted(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneByIDAndUserIDAndNotDeleted", reflect.TypeOf((*MockAgentRepository)(nil).FindOneByIDAndUserIDAndNotDeleted), arg0, arg1, arg2)
}

// FindOneByTokenAndNotDeleted mocks base method.
func (m *MockAgentRepository) FindOneByTokenAndNotDeleted(arg0 context.Context, arg1 string) (*entity.Agent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneByTokenAndNotDeleted", arg0, arg1)
	ret0, _ := ret[0].(*entity.Agent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneByTokenAndNotDeleted indicates an expected call of FindOneByTokenAndNotDeleted.
func (mr *MockAgentRepositoryMockRecorder) FindOneByTokenAndNotDeleted(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneByTokenAndNotDeleted", reflect.TypeOf((*MockAgentRepository)(nil).FindOneByTokenAndNotDeleted), arg0, arg1)
}

// Update mocks base method.
func (m *MockAgentRepository) Update(arg0 context.Context, arg1 *entity.Agent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockAgentRepositoryMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAgentRepository)(nil).Update), arg0, arg1)
}
