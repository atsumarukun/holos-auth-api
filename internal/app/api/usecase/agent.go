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
	Gets(context.Context, uuid.UUID) ([]*dto.AgentDTO, apierr.ApiError)
	UpdatePolicies(context.Context, uuid.UUID, uuid.UUID, []uuid.UUID) ([]*dto.PolicyDTO, apierr.ApiError)
	GetPolicies(context.Context, uuid.UUID, uuid.UUID) ([]*dto.PolicyDTO, apierr.ApiError)
}

type agentUsecase struct {
	transactionObject domain.TransactionObject
	agentRepository   repository.AgentRepository
	policyrepository  repository.PolicyRepository
}

func NewAgentUsecase(transactionObject domain.TransactionObject, agentRepository repository.AgentRepository, policyrepository repository.PolicyRepository) AgentUsecase {
	return &agentUsecase{
		transactionObject: transactionObject,
		agentRepository:   agentRepository,
		policyrepository:  policyrepository,
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

func (u *agentUsecase) Gets(ctx context.Context, userID uuid.UUID) ([]*dto.AgentDTO, apierr.ApiError) {
	agents, err := u.agentRepository.FindByUserIDAndNotDeleted(ctx, userID)
	if err != nil {
		return nil, err
	}

	return u.convertToDTOs(agents), nil
}

func (u *agentUsecase) UpdatePolicies(ctx context.Context, id uuid.UUID, userID uuid.UUID, policyIDs []uuid.UUID) ([]*dto.PolicyDTO, apierr.ApiError) {
	policies := make([]*entity.Policy, len(policyIDs))

	if err := u.transactionObject.Transaction(ctx, func(ctx context.Context) apierr.ApiError {
		agent, err := u.agentRepository.FindOneByIDAndUserIDAndNotDeleted(ctx, id, userID)
		if err != nil {
			return err
		}

		policies, err = u.policyrepository.FindByIDsAndUserIDAndNotDeleted(ctx, policyIDs, userID)
		if err != nil {
			return err
		}

		return u.agentRepository.UpdatePolicies(ctx, agent.ID, policies)
	}); err != nil {
		return nil, err
	}

	dtos := make([]*dto.PolicyDTO, len(policies))
	for i, policy := range policies {
		dtos[i] = dto.NewPolicyDTO(policy.ID, policy.UserID, policy.Name, policy.Service, policy.Path, policy.Methods, policy.CreatedAt, policy.UpdatedAt)
	}

	return dtos, nil
}

func (u *agentUsecase) GetPolicies(ctx context.Context, id uuid.UUID, userID uuid.UUID) ([]*dto.PolicyDTO, apierr.ApiError) {
	policies, err := u.agentRepository.GetPolicies(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	dtos := make([]*dto.PolicyDTO, len(policies))
	for i, policy := range policies {
		dtos[i] = dto.NewPolicyDTO(policy.ID, policy.UserID, policy.Name, policy.Service, policy.Path, policy.Methods, policy.CreatedAt, policy.UpdatedAt)
	}
	return dtos, nil
}

func (u *agentUsecase) convertToDTO(agent *entity.Agent) *dto.AgentDTO {
	return dto.NewAgentDTO(agent.ID, agent.UserID, agent.Name, agent.CreatedAt, agent.UpdatedAt)
}

func (u *agentUsecase) convertToDTOs(agents []*entity.Agent) []*dto.AgentDTO {
	dtos := make([]*dto.AgentDTO, len(agents))
	for i, agent := range agents {
		dtos[i] = u.convertToDTO(agent)
	}
	return dtos
}
