//go:generate mockgen -source=$GOFILE -destination=../../../../../test/mock/domain/service/$GOFILE
package service

import (
	"context"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/repository"
	"holos-auth-api/internal/pkg/apierr"
)

type AgentService interface {
	Exists(context.Context, *entity.Agent) (bool, apierr.ApiError)
}

type agentService struct {
	agentRepository repository.AgentRepository
}

func NewAgentService(agentRepository repository.AgentRepository) AgentService {
	return &agentService{
		agentRepository: agentRepository,
	}
}

func (as *agentService) Exists(ctx context.Context, agent *entity.Agent) (bool, apierr.ApiError) {
	agent, err := as.agentRepository.FindOneByUserIDAndName(ctx, agent.UserID, agent.Name)
	if err != nil {
		return false, err
	}
	return agent != nil, err
}
