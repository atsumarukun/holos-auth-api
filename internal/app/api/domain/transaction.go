package domain

import (
	"context"
)

type TransactionObject interface {
	Transaction(context.Context, func(context.Context) error) error
}
