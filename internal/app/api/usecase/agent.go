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
	ErrAgentAlreadyExists = status.Error(http.StatusBadRequest, "agent already exists")
	ErrAgentNotFound      = status.Error(http.StatusNotFound, "agent not found")
	ErrAgentTokenNotFound = status.Error(http.StatusNotFound, "agent token not found")
)

type AgentUsecase interface {
	Create(context.Context, uuid.UUID, string) (*dto.AgentDTO, error)
	Update(context.Context, uuid.UUID, uuid.UUID, string) (*dto.AgentDTO, error)
	Delete(context.Context, uuid.UUID, uuid.UUID) error
	Gets(context.Context, uuid.UUID) ([]*dto.AgentDTO, error)
	UpdatePolicies(context.Context, uuid.UUID, uuid.UUID, []uuid.UUID) ([]*dto.PolicyDTO, error)
	GetPolicies(context.Context, uuid.UUID, uuid.UUID) ([]*dto.PolicyDTO, error)
	GenerateToken(context.Context, uuid.UUID, uuid.UUID) (string, error)
	DeleteToken(context.Context, uuid.UUID, uuid.UUID) error
}

type agentUsecase struct {
	transactionObject    domain.TransactionObject
	agentRepository      repository.AgentRepository
	agentTokenRepository repository.AgentTokenRepository
	policyRepository     repository.PolicyRepository
	agentService         service.AgentService
}

func NewAgentUsecase(
	transactionObject domain.TransactionObject,
	agentRepository repository.AgentRepository,
	agentTokenRepository repository.AgentTokenRepository,
	policyRepository repository.PolicyRepository,
	agentService service.AgentService,
) AgentUsecase {
	return &agentUsecase{
		transactionObject:    transactionObject,
		agentRepository:      agentRepository,
		agentTokenRepository: agentTokenRepository,
		policyRepository:     policyRepository,
		agentService:         agentService,
	}
}

func (u *agentUsecase) Create(ctx context.Context, userID uuid.UUID, name string) (*dto.AgentDTO, error) {
	agent, err := entity.NewAgent(userID, name)
	if err != nil {
		return nil, err
	}

	if err := u.agentRepository.Create(ctx, agent); err != nil {
		return nil, err
	}

	return mapper.ToAgentDTO(agent), nil
}

func (u *agentUsecase) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, name string) (*dto.AgentDTO, error) {
	var agent *entity.Agent

	if err := u.transactionObject.Transaction(ctx, func(ctx context.Context) error {
		var err error
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

	return mapper.ToAgentDTO(agent), nil
}

func (u *agentUsecase) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return u.transactionObject.Transaction(ctx, func(ctx context.Context) error {
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

func (u *agentUsecase) Gets(ctx context.Context, userID uuid.UUID) ([]*dto.AgentDTO, error) {
	agents, err := u.agentRepository.FindByUserIDAndNotDeleted(ctx, userID)
	if err != nil {
		return nil, err
	}

	return mapper.ToAgentDTOs(agents), nil
}

func (u *agentUsecase) UpdatePolicies(ctx context.Context, id uuid.UUID, userID uuid.UUID, policyIDs []uuid.UUID) ([]*dto.PolicyDTO, error) {
	policies := make([]*entity.Policy, len(policyIDs))

	if err := u.transactionObject.Transaction(ctx, func(ctx context.Context) error {
		agent, err := u.agentRepository.FindOneByIDAndUserIDAndNotDeleted(ctx, id, userID)
		if err != nil {
			return err
		}
		if agent == nil {
			return ErrAgentNotFound
		}

		policies, err = u.policyRepository.FindByIDsAndUserIDAndNotDeleted(ctx, policyIDs, userID)
		if err != nil {
			return err
		}

		agent.SetPolicies(policies)

		return u.agentRepository.Update(ctx, agent)
	}); err != nil {
		return nil, err
	}

	return mapper.ToPolicyDTOs(policies), nil
}

func (u *agentUsecase) GetPolicies(ctx context.Context, id uuid.UUID, userID uuid.UUID) ([]*dto.PolicyDTO, error) {
	policies := []*entity.Policy{}

	if err := u.transactionObject.Transaction(ctx, func(ctx context.Context) error {
		agent, err := u.agentRepository.FindOneByIDAndUserIDAndNotDeleted(ctx, id, userID)
		if err != nil {
			return err
		}
		if agent == nil {
			return ErrAgentNotFound
		}

		policies, err = u.agentService.GetPolicies(ctx, agent)
		return err
	}); err != nil {
		return nil, err
	}

	return mapper.ToPolicyDTOs(policies), nil
}

func (u *agentUsecase) GenerateToken(ctx context.Context, id uuid.UUID, userID uuid.UUID) (string, error) {
	var agentToken *entity.AgentToken

	if err := u.transactionObject.Transaction(ctx, func(ctx context.Context) error {
		agent, err := u.agentRepository.FindOneByIDAndUserIDAndNotDeleted(ctx, id, userID)
		if err != nil {
			return err
		}
		if agent == nil {
			return ErrAgentNotFound
		}

		agentToken, err = entity.NewAgentToken(agent.ID)
		if err != nil {
			return err
		}

		return u.agentTokenRepository.Save(ctx, agentToken)
	}); err != nil {
		return "", err
	}

	return agentToken.Token, nil
}

func (u *agentUsecase) DeleteToken(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return u.transactionObject.Transaction(ctx, func(ctx context.Context) error {
		agentToken, err := u.agentTokenRepository.FindOneByAgentIDAndUserID(ctx, id, userID)
		if err != nil {
			return err
		}
		if agentToken == nil {
			return ErrAgentTokenNotFound
		}

		return u.agentTokenRepository.Delete(ctx, agentToken)
	})
}
