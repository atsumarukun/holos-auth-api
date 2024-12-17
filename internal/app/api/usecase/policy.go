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
	ErrPolicyNotFound = apierr.NewApiError(http.StatusNotFound, "policy not found")
)

type PolicyUsecase interface {
	Create(context.Context, uuid.UUID, string, string, string, []string) (*dto.PolicyDTO, error)
	Update(context.Context, uuid.UUID, uuid.UUID, string, string, string, []string) (*dto.PolicyDTO, error)
	Delete(context.Context, uuid.UUID, uuid.UUID) error
	Gets(context.Context, uuid.UUID) ([]*dto.PolicyDTO, error)
	UpdateAgents(context.Context, uuid.UUID, uuid.UUID, []uuid.UUID) ([]*dto.AgentDTO, error)
	GetAgents(context.Context, uuid.UUID, uuid.UUID) ([]*dto.AgentDTO, error)
}

type policyUsecase struct {
	transactionObject domain.TransactionObject
	policyRepository  repository.PolicyRepository
	agentRepository   repository.AgentRepository
}

func NewPolicyUsecase(transactionObject domain.TransactionObject, policyRepository repository.PolicyRepository, agentRepository repository.AgentRepository) PolicyUsecase {
	return &policyUsecase{
		transactionObject: transactionObject,
		policyRepository:  policyRepository,
		agentRepository:   agentRepository,
	}
}

func (u *policyUsecase) Create(ctx context.Context, userID uuid.UUID, name string, service string, path string, methods []string) (*dto.PolicyDTO, error) {
	policy, err := entity.NewPolicy(userID, name, service, path, methods)
	if err != nil {
		return nil, err
	}

	if err := u.policyRepository.Create(ctx, policy); err != nil {
		return nil, err
	}

	return u.convertToDTO(policy), nil
}

func (u *policyUsecase) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, name string, service string, path string, methods []string) (*dto.PolicyDTO, error) {
	var policy *entity.Policy

	if err := u.transactionObject.Transaction(ctx, func(ctx context.Context) error {
		var err error
		policy, err = u.policyRepository.FindOneByIDAndUserIDAndNotDeleted(ctx, id, userID)
		if err != nil {
			return err
		}
		if policy == nil {
			return ErrPolicyNotFound
		}

		if err := policy.SetName(name); err != nil {
			return err
		}
		if err := policy.SetService(service); err != nil {
			return err
		}
		if err := policy.SetPath(path); err != nil {
			return err
		}
		if err := policy.SetMethods(methods); err != nil {
			return err
		}

		return u.policyRepository.Update(ctx, policy)
	}); err != nil {
		return nil, err
	}

	return u.convertToDTO(policy), nil
}

func (u *policyUsecase) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return u.transactionObject.Transaction(ctx, func(ctx context.Context) error {
		policy, err := u.policyRepository.FindOneByIDAndUserIDAndNotDeleted(ctx, id, userID)
		if err != nil {
			return err
		}
		if policy == nil {
			return ErrPolicyNotFound
		}

		return u.policyRepository.Delete(ctx, policy)
	})
}

func (u *policyUsecase) Gets(ctx context.Context, userID uuid.UUID) ([]*dto.PolicyDTO, error) {
	policies, err := u.policyRepository.FindByUserIDAndNotDeleted(ctx, userID)
	if err != nil {
		return nil, err
	}

	return u.convertToDTOs(policies), nil
}

func (u *policyUsecase) UpdateAgents(ctx context.Context, id uuid.UUID, userID uuid.UUID, agentIDs []uuid.UUID) ([]*dto.AgentDTO, error) {
	agents := make([]*entity.Agent, len(agentIDs))

	if err := u.transactionObject.Transaction(ctx, func(ctx context.Context) error {
		policy, err := u.policyRepository.FindOneByIDAndUserIDAndNotDeleted(ctx, id, userID)
		if err != nil {
			return err
		}
		if policy == nil {
			return ErrPolicyNotFound
		}

		agents, err = u.agentRepository.FindByIDsAndUserIDAndNotDeleted(ctx, agentIDs, userID)
		if err != nil {
			return err
		}

		policy.SetAgents(agents)

		return u.policyRepository.Update(ctx, policy)
	}); err != nil {
		return nil, err
	}

	dtos := make([]*dto.AgentDTO, len(agents))
	for i, agent := range agents {
		dtos[i] = dto.NewAgentDTO(agent.ID, agent.UserID, agent.Name, agent.CreatedAt, agent.UpdatedAt)
	}
	return dtos, nil
}

func (u *policyUsecase) GetAgents(ctx context.Context, id uuid.UUID, userID uuid.UUID) ([]*dto.AgentDTO, error) {
	policy, err := u.policyRepository.FindOneByIDAndUserIDAndNotDeleted(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return nil, ErrPolicyNotFound
	}

	agents, err := u.agentRepository.FindByIDsAndUserIDAndNotDeleted(ctx, policy.Agents, userID)
	if err != nil {
		return nil, err
	}

	dtos := make([]*dto.AgentDTO, len(agents))
	for i, agent := range agents {
		dtos[i] = dto.NewAgentDTO(agent.ID, agent.UserID, agent.Name, agent.CreatedAt, agent.UpdatedAt)
	}
	return dtos, nil
}

func (u *policyUsecase) convertToDTO(policy *entity.Policy) *dto.PolicyDTO {
	return dto.NewPolicyDTO(policy.ID, policy.UserID, policy.Name, policy.Service, policy.Path, policy.Methods, policy.CreatedAt, policy.UpdatedAt)
}

func (u *policyUsecase) convertToDTOs(policies []*entity.Policy) []*dto.PolicyDTO {
	dtos := make([]*dto.PolicyDTO, len(policies))
	for i, policy := range policies {
		dtos[i] = u.convertToDTO(policy)
	}
	return dtos
}
