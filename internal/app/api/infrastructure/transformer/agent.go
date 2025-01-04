package transformer

import (
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/infrastructure/model"
	"strings"

	"github.com/google/uuid"
)

func ToAgentModel(agent *entity.Agent) *model.AgentModel {
	var policies string
	if len(agent.Policies) != 0 {
		policyIDs := make([]string, len(agent.Policies))
		for i, policy := range agent.Policies {
			policyIDs[i] = policy.String()
		}
		policies = strings.Join(policyIDs, ",")
	}

	return &model.AgentModel{
		ID:        agent.ID,
		UserID:    agent.UserID,
		Name:      agent.Name,
		Policies:  &policies,
		CreatedAt: agent.CreatedAt,
		UpdatedAt: agent.UpdatedAt,
	}
}

func ToAgentEntity(agent *model.AgentModel) (*entity.Agent, error) {
	policies := []uuid.UUID{}
	if agent.Policies != nil {
		for _, policyID := range strings.Split(*agent.Policies, ",") {
			policy, err := uuid.Parse(policyID)
			if err != nil {
				return nil, err
			}
			policies = append(policies, policy)
		}
	}

	return entity.RestoreAgent(
		agent.ID,
		agent.UserID,
		agent.Name,
		policies,
		agent.CreatedAt,
		agent.UpdatedAt,
	), nil
}

func ToAgentEntities(agents []*model.AgentModel) ([]*entity.Agent, error) {
	entities := make([]*entity.Agent, len(agents))
	var err error
	for i, agent := range agents {
		entities[i], err = ToAgentEntity(agent)
		if err != nil {
			return nil, err
		}
	}
	return entities, nil
}
