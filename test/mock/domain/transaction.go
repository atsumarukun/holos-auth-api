// Code generated by MockGen. DO NOT EDIT.
// Source: transaction.go

// Package mock_domain is a generated GoMock package.
package mock_domain

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockTransactionObject is a mock of TransactionObject interface.
type MockTransactionObject struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionObjectMockRecorder
}

// MockTransactionObjectMockRecorder is the mock recorder for MockTransactionObject.
type MockTransactionObjectMockRecorder struct {
	mock *MockTransactionObject
}

// NewMockTransactionObject creates a new mock instance.
func NewMockTransactionObject(ctrl *gomock.Controller) *MockTransactionObject {
	mock := &MockTransactionObject{ctrl: ctrl}
	mock.recorder = &MockTransactionObjectMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionObject) EXPECT() *MockTransactionObjectMockRecorder {
	return m.recorder
}

// Transaction mocks base method.
func (m *MockTransactionObject) Transaction(arg0 context.Context, arg1 func(context.Context) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transaction", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Transaction indicates an expected call of Transaction.
func (mr *MockTransactionObjectMockRecorder) Transaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transaction", reflect.TypeOf((*MockTransactionObject)(nil).Transaction), arg0, arg1)
}