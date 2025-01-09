//go:generate mockgen -source=$GOFILE -destination=../../../../../test/mock/domain/repository/$GOFILE
package repository

import (
	"context"
	"holos-auth-api/internal/app/api/domain/entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(context.Context, *entity.User) error
	Update(context.Context, *entity.User) error
	Delete(context.Context, *entity.User) error
	FindOneByIDAndNotDeleted(context.Context, uuid.UUID) (*entity.User, error)
	FindOneByName(context.Context, string) (*entity.User, error)
}
