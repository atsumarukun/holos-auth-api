//go:generate mockgen -source=$GOFILE -destination=../../../../../test/mock/domain/repository/$GOFILE
package repository

import (
	"context"
	"holos-auth-api/internal/app/api/domain/entity"

	"github.com/google/uuid"
)

type AgentTokenRepository interface {
	Save(context.Context, *entity.AgentToken) error
	Delete(context.Context, *entity.AgentToken) error
	FindOneByAgentIDAndUserID(context.Context, uuid.UUID, uuid.UUID) (*entity.AgentToken, error)
}
