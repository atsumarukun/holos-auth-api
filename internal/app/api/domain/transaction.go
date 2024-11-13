package domain

import (
	"context"
	"holos-auth-api/internal/pkg/apierr"
)

type TransactionObject interface {
	Transaction(context.Context, func(context.Context) apierr.ApiError) apierr.ApiError
}
