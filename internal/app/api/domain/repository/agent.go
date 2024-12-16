//go:generate mockgen -source=$GOFILE -destination=../../../../../test/mock/domain/repository/$GOFILE
package repository

import (
	"context"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/pkg/apierr"

	"github.com/google/uuid"
)

type AgentRepository interface {
	Create(context.Context, *entity.Agent) apierr.ApiError
	Update(context.Context, *entity.Agent) apierr.ApiError
	Delete(context.Context, *entity.Agent) apierr.ApiError
	FindOneByIDAndUserIDAndNotDeleted(context.Context, uuid.UUID, uuid.UUID) (*entity.Agent, apierr.ApiError)
	FindByUserIDAndNotDeleted(context.Context, uuid.UUID) ([]*entity.Agent, apierr.ApiError)
	FindByIDsAndUserIDAndNotDeleted(context.Context, []uuid.UUID, uuid.UUID) ([]*entity.Agent, apierr.ApiError)
}
