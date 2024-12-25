//go:generate mockgen -source=$GOFILE -destination=../../../../test/mock/domain/$GOFILE
package domain

import (
	"context"
)

type TransactionObject interface {
	Transaction(context.Context, func(context.Context) error) error
}
