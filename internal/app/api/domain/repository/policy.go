//go:generate mockgen -source=$GOFILE -destination=../../../../../test/mock/domain/repository/$GOFILE
package repository

import (
	"context"
	"holos-auth-api/internal/app/api/domain/entity"

	"github.com/google/uuid"
)

type PolicyRepository interface {
	Create(context.Context, *entity.Policy) error
	Update(context.Context, *entity.Policy) error
	Delete(context.Context, *entity.Policy) error
	FindOneByIDAndUserIDAndNotDeleted(context.Context, uuid.UUID, uuid.UUID) (*entity.Policy, error)
	FindByNamePrefixAndUserIDAndNotDeleted(context.Context, string, uuid.UUID) ([]*entity.Policy, error)
	FindByIDsAndUserIDAndNotDeleted(context.Context, []uuid.UUID, uuid.UUID) ([]*entity.Policy, error)
}
