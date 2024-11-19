package domain

import (
	"context"
	"holos-auth-api/internal/app/api/pkg/apierr"
)

type TransactionObject interface {
	Transaction(context.Context, func(context.Context) apierr.ApiError) apierr.ApiError
}
