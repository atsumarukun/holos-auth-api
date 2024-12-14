//go:generate mockgen -source=$GOFILE -destination=../../../../../test/mock/domain/repository/$GOFILE
package repository

import (
	"context"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/pkg/apierr"

	"github.com/google/uuid"
)

type PolicyRepository interface {
	Create(context.Context, *entity.Policy) apierr.ApiError
	Update(context.Context, *entity.Policy) apierr.ApiError
	Delete(context.Context, *entity.Policy) apierr.ApiError
	FindOneByIDAndUserIDAndNotDeleted(context.Context, uuid.UUID, uuid.UUID) (*entity.Policy, apierr.ApiError)
	FindByUserIDAndNotDeleted(context.Context, uuid.UUID) ([]*entity.Policy, apierr.ApiError)
	FindByIDsAndUserIDAndNotDeleted(context.Context, []uuid.UUID, uuid.UUID) ([]*entity.Policy, apierr.ApiError)
	UpdateAgents(context.Context, uuid.UUID, []*entity.Agent) apierr.ApiError
	GetAgents(context.Context, uuid.UUID, uuid.UUID) ([]*entity.Agent, apierr.ApiError)
}
