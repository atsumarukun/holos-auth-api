package test

import (
	"context"
	"holos-auth-api/internal/app/api/domain"
	"holos-auth-api/internal/pkg/apierr"
)

type testTransactionObject struct{}

func NewTestTransactionObject() domain.TransactionObject {
	return &testTransactionObject{}
}

func (tto *testTransactionObject) Transaction(ctx context.Context, fn func(context.Context) apierr.ApiError) apierr.ApiError {
	return fn(ctx)
}
