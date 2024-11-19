//go:generate mockgen -source=$GOFILE -destination=../../../../../test/mock/domain/repository/$GOFILE
package repository

import (
	"context"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/pkg/apierr"

	"github.com/google/uuid"
)

type AgentRepository interface {
	Create(context.Context, *entity.Agent) apierr.ApiError
	Update(context.Context, *entity.Agent) apierr.ApiError
	Delete(context.Context, *entity.Agent) apierr.ApiError
	FindOneByIDAndUserIDAndNotDeleted(context.Context, uuid.UUID, uuid.UUID) (*entity.Agent, apierr.ApiError)
	FindOneByUserIDAndName(context.Context, uuid.UUID, string) (*entity.Agent, apierr.ApiError)
}
