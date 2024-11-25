//go:generate mockgen -source=$GOFILE -destination=../../../../test/mock/usecase/$GOFILE
package usecase

import (
	"context"
	"holos-auth-api/internal/app/api/domain"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/repository"
	"holos-auth-api/internal/app/api/pkg/apierr"
	"holos-auth-api/internal/app/api/usecase/dto"
	"net/http"

	"github.com/google/uuid"
)

var (
	ErrAgentAlreadyExists = apierr.NewApiError(http.StatusBadRequest, "agent already exists")
	ErrAgentNotFound      = apierr.NewApiError(http.StatusNotFound, "agent not found")
)

type AgentUsecase interface {
	Create(context.Context, uuid.UUID, string) (*dto.AgentDTO, apierr.ApiError)
	Update(context.Context, uuid.UUID, uuid.UUID, string) (*dto.AgentDTO, apierr.ApiError)
	Delete(context.Context, uuid.UUID, uuid.UUID) apierr.ApiError
}

type agentUsecase struct {
	transactionObject domain.TransactionObject
	agentRepository   repository.AgentRepository
}

func NewAgentUsecase(transactionObject domain.TransactionObject, agentRepository repository.AgentRepository) AgentUsecase {
	return &agentUsecase{
		transactionObject: transactionObject,
		agentRepository:   agentRepository,
	}
}

func (u *agentUsecase) Create(ctx context.Context, userID uuid.UUID, name string) (*dto.AgentDTO, apierr.ApiError) {
	agent, err := entity.NewAgent(userID, name)
	if err != nil {
		return nil, err
	}

	if err := u.agentRepository.Create(ctx, agent); err != nil {
		return nil, err
	}

	return dto.NewAgentDTO(agent.ID, agent.UserID, agent.Name, agent.CreatedAt, agent.UpdatedAt), nil
}

func (u *agentUsecase) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, name string) (*dto.AgentDTO, apierr.ApiError) {
	var agent *entity.Agent

	if err := u.transactionObject.Transaction(ctx, func(ctx context.Context) apierr.ApiError {
		var err apierr.ApiError
		agent, err = u.agentRepository.FindOneByIDAndUserIDAndNotDeleted(ctx, id, userID)
		if err != nil {
			return err
		}
		if agent == nil {
			return ErrAgentNotFound
		}

		if err := agent.SetName(name); err != nil {
			return err
		}

		return u.agentRepository.Update(ctx, agent)
	}); err != nil {
		return nil, err
	}

	return dto.NewAgentDTO(agent.ID, agent.UserID, agent.Name, agent.CreatedAt, agent.UpdatedAt), nil
}

func (u *agentUsecase) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) apierr.ApiError {
	return u.transactionObject.Transaction(ctx, func(ctx context.Context) apierr.ApiError {
		agent, err := u.agentRepository.FindOneByIDAndUserIDAndNotDeleted(ctx, id, userID)
		if err != nil {
			return err
		}
		if agent == nil {
			return ErrAgentNotFound
		}

		return u.agentRepository.Delete(ctx, agent)
	})
}
