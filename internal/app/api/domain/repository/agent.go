//go:generate mockgen -source=$GOFILE -destination=../../../../../test/mock/domain/repository/$GOFILE
package repository

import (
	"context"
	"holos-auth-api/internal/app/api/domain/entity"

	"github.com/google/uuid"
)

type AgentRepository interface {
	Create(context.Context, *entity.Agent) error
	Update(context.Context, *entity.Agent) error
	Delete(context.Context, *entity.Agent) error
	FindOneByIDAndUserIDAndNotDeleted(context.Context, uuid.UUID, uuid.UUID) (*entity.Agent, error)
	FindOneByTokenAndNotDeleted(context.Context, string) (*entity.Agent, error)
	FindByUserIDAndNotDeleted(context.Context, uuid.UUID) ([]*entity.Agent, error)
	FindByIDsAndUserIDAndNotDeleted(context.Context, []uuid.UUID, uuid.UUID) ([]*entity.Agent, error)
}
