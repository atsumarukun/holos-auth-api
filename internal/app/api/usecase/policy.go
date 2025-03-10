//go:generate mockgen -source=$GOFILE -destination=../../../../test/mock/usecase/$GOFILE
package usecase

import (
	"context"
	"holos-auth-api/internal/app/api/domain"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/repository"
	"holos-auth-api/internal/app/api/domain/service"
	"holos-auth-api/internal/app/api/pkg/status"
	"holos-auth-api/internal/app/api/usecase/dto"
	"holos-auth-api/internal/app/api/usecase/mapper"
	"net/http"

	"github.com/google/uuid"
)

var (
	ErrPolicyNotFound = status.Error(http.StatusNotFound, "policy not found")
)

type PolicyUsecase interface {
	Create(context.Context, uuid.UUID, string, string, string, string, []string) (*dto.PolicyDTO, error)
	Update(context.Context, uuid.UUID, uuid.UUID, string, string, string, string, []string) (*dto.PolicyDTO, error)
	Delete(context.Context, uuid.UUID, uuid.UUID) error
	Get(context.Context, uuid.UUID, uuid.UUID) (*dto.PolicyDTO, error)
	Gets(context.Context, string, uuid.UUID) ([]*dto.PolicyDTO, error)
	UpdateAgents(context.Context, uuid.UUID, uuid.UUID, []uuid.UUID) ([]*dto.AgentDTO, error)
	GetAgents(context.Context, uuid.UUID, uuid.UUID, string) ([]*dto.AgentDTO, error)
}

type policyUsecase struct {
	transactionObject domain.TransactionObject
	policyRepository  repository.PolicyRepository
	agentRepository   repository.AgentRepository
	policyService     service.PolicyService
}

func NewPolicyUsecase(
	transactionObject domain.TransactionObject,
	policyRepository repository.PolicyRepository,
	agentRepository repository.AgentRepository,
	policyService service.PolicyService,
) PolicyUsecase {
	return &policyUsecase{
		transactionObject: transactionObject,
		policyRepository:  policyRepository,
		agentRepository:   agentRepository,
		policyService:     policyService,
	}
}

func (u *policyUsecase) Create(ctx context.Context, userID uuid.UUID, name string, effect string, service string, path string, methods []string) (*dto.PolicyDTO, error) {
	policy, err := entity.NewPolicy(userID, name, effect, service, path, methods)
	if err != nil {
		return nil, err
	}

	if err := u.policyRepository.Create(ctx, policy); err != nil {
		return nil, err
	}

	return mapper.ToPolicyDTO(policy), nil
}

func (u *policyUsecase) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, name string, effect string, service string, path string, methods []string) (*dto.PolicyDTO, error) {
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
		if err := policy.SetEffect(effect); err != nil {
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

	return mapper.ToPolicyDTO(policy), nil
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

func (u *policyUsecase) Get(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.PolicyDTO, error) {
	policy, err := u.policyRepository.FindOneByIDAndUserIDAndNotDeleted(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return nil, ErrPolicyNotFound
	}

	return mapper.ToPolicyDTO(policy), nil
}

func (u *policyUsecase) Gets(ctx context.Context, keyword string, userID uuid.UUID) ([]*dto.PolicyDTO, error) {
	policies, err := u.policyRepository.FindByNamePrefixAndUserIDAndNotDeleted(ctx, keyword, userID)
	if err != nil {
		return nil, err
	}

	return mapper.ToPolicyDTOs(policies), nil
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

	return mapper.ToAgentDTOs(agents), nil
}

func (u *policyUsecase) GetAgents(ctx context.Context, id uuid.UUID, userID uuid.UUID, keyword string) ([]*dto.AgentDTO, error) {
	agents := []*entity.Agent{}

	if err := u.transactionObject.Transaction(ctx, func(ctx context.Context) error {
		policy, err := u.policyRepository.FindOneByIDAndUserIDAndNotDeleted(ctx, id, userID)
		if err != nil {
			return err
		}
		if policy == nil {
			return ErrPolicyNotFound
		}

		agents, err = u.policyService.GetAgents(ctx, policy, keyword)
		return err
	}); err != nil {
		return nil, err
	}

	return mapper.ToAgentDTOs(agents), nil
}
