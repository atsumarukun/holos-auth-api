//go:generate mockgen -source=$GOFILE -destination=../../../../../test/mock/domain/service/$GOFILE
package service

import (
	"context"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/repository"
	"holos-auth-api/internal/app/api/pkg/status"
	"net/http"
	"regexp"
	"slices"
	"sort"
)

var (
	ErrRequiredAgent = status.Error(http.StatusInternalServerError, "agent is required")
)

type AgentService interface {
	GetPolicies(context.Context, *entity.Agent) ([]*entity.Policy, error)
	HasPermission(context.Context, *entity.Agent, string, string, string) (bool, error)
}

type agentService struct {
	policyRepository repository.PolicyRepository
}

func NewAgentService(policyRepository repository.PolicyRepository) AgentService {
	return &agentService{
		policyRepository: policyRepository,
	}
}

func (s *agentService) GetPolicies(ctx context.Context, agent *entity.Agent) ([]*entity.Policy, error) {
	if agent == nil {
		return nil, ErrRequiredAgent
	}

	return s.policyRepository.FindByIDsAndUserIDAndNotDeleted(ctx, agent.Policies, agent.UserID)
}

func (s *agentService) HasPermission(ctx context.Context, agent *entity.Agent, service string, path string, method string) (bool, error) {
	policies, err := s.GetPolicies(ctx, agent)
	if err != nil {
		return false, err
	}
	if len(policies) == 0 {
		return false, nil
	}

	sort.Slice(policies, func(i, j int) bool {
		return policies[j].Path < policies[i].Path
	})

	for _, policy := range policies {
		if service != policy.Service || !slices.Contains(policy.Methods, method) {
			continue
		}

		pattern, err := regexp.Compile(`:[^/]+`)
		if err != nil {
			return false, err
		}
		matched, err := regexp.MatchString("^"+pattern.ReplaceAllString(policy.Path, `[^/]+`), path)
		if err != nil {
			return false, err
		}
		if matched {
			if policy.Effect == "ALLOW" {
				return true, nil
			}
			break
		}
	}

	return false, nil
}
