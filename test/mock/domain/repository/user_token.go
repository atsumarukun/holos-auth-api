// Code generated by MockGen. DO NOT EDIT.
// Source: user_token.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	entity "holos-auth-api/internal/app/api/domain/entity"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserTokenRepository is a mock of UserTokenRepository interface.
type MockUserTokenRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserTokenRepositoryMockRecorder
}

// MockUserTokenRepositoryMockRecorder is the mock recorder for MockUserTokenRepository.
type MockUserTokenRepositoryMockRecorder struct {
	mock *MockUserTokenRepository
}

// NewMockUserTokenRepository creates a new mock instance.
func NewMockUserTokenRepository(ctrl *gomock.Controller) *MockUserTokenRepository {
	mock := &MockUserTokenRepository{ctrl: ctrl}
	mock.recorder = &MockUserTokenRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserTokenRepository) EXPECT() *MockUserTokenRepositoryMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockUserTokenRepository) Delete(arg0 context.Context, arg1 *entity.UserToken) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUserTokenRepositoryMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUserTokenRepository)(nil).Delete), arg0, arg1)
}

// FindOneByTokenAndNotExpired mocks base method.
func (m *MockUserTokenRepository) FindOneByTokenAndNotExpired(arg0 context.Context, arg1 string) (*entity.UserToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneByTokenAndNotExpired", arg0, arg1)
	ret0, _ := ret[0].(*entity.UserToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneByTokenAndNotExpired indicates an expected call of FindOneByTokenAndNotExpired.
func (mr *MockUserTokenRepositoryMockRecorder) FindOneByTokenAndNotExpired(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneByTokenAndNotExpired", reflect.TypeOf((*MockUserTokenRepository)(nil).FindOneByTokenAndNotExpired), arg0, arg1)
}

// Save mocks base method.
func (m *MockUserTokenRepository) Save(arg0 context.Context, arg1 *entity.UserToken) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockUserTokenRepositoryMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockUserTokenRepository)(nil).Save), arg0, arg1)
}
