package test

import (
	"context"
	"holos-auth-api/internal/app/api/domain"
)

type testTransactionObject struct{}

func NewTestTransactionObject() domain.TransactionObject {
	return &testTransactionObject{}
}

func (tto *testTransactionObject) Transaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}
