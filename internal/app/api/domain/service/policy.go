//go:generate mockgen -source=$GOFILE -destination=../../../../../test/mock/domain/service/$GOFILE
package service

import (
	"context"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/repository"
	"holos-auth-api/internal/app/api/pkg/status"
	"net/http"
)

var (
	ErrRequiredPolicy = status.Error(http.StatusInternalServerError, "policy is required")
)

type PolicyService interface {
	GetAgents(context.Context, *entity.Policy, string) ([]*entity.Agent, error)
}

type policyService struct {
	agentRepository repository.AgentRepository
}

func NewPolicyService(agentRepository repository.AgentRepository) PolicyService {
	return &policyService{
		agentRepository: agentRepository,
	}
}

func (s *policyService) GetAgents(ctx context.Context, polisy *entity.Policy, keyword string) ([]*entity.Agent, error) {
	if polisy == nil {
		return nil, ErrRequiredPolicy
	}

	return s.agentRepository.FindByIDsAndNamePrefixAndUserIDAndNotDeleted(ctx, polisy.Agents, keyword, polisy.UserID)
}
