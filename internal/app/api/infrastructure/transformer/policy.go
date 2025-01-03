package transformer

import (
	"encoding/json"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/infrastructure/model"
	"holos-auth-api/internal/app/api/pkg/status"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func ToPolicyModel(policy *entity.Policy) (*model.PolicyModel, error) {
	methods, err := json.Marshal(policy.Methods)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	var agents string
	if len(policy.Agents) != 0 {
		agentIDs := make([]string, len(policy.Agents))
		for i, agent := range policy.Agents {
			agentIDs[i] = agent.String()
		}
		agents = strings.Join(agentIDs, ",")
	}

	return &model.PolicyModel{
		ID:        policy.ID,
		UserID:    policy.UserID,
		Name:      policy.Name,
		Service:   policy.Service,
		Path:      policy.Path,
		Methods:   methods,
		Agents:    &agents,
		CreatedAt: policy.CreatedAt,
		UpdatedAt: policy.UpdatedAt,
	}, nil
}

func ToPolicyEntity(policy *model.PolicyModel) (*entity.Policy, error) {
	var methods []string
	if err := json.Unmarshal(policy.Methods, &methods); err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	agents := []uuid.UUID{}
	if policy.Agents != nil {
		for _, agentID := range strings.Split(*policy.Agents, ",") {
			agent, err := uuid.Parse(agentID)
			if err != nil {
				return nil, status.Error(http.StatusInternalServerError, err.Error())
			}
			agents = append(agents, agent)
		}
	}

	return entity.RestorePolicy(
		policy.ID,
		policy.UserID,
		policy.Name,
		policy.Service,
		policy.Path,
		methods,
		agents,
		policy.CreatedAt,
		policy.UpdatedAt,
	), nil
}

func ToPolicyEntities(policies []*model.PolicyModel) ([]*entity.Policy, error) {
	entities := make([]*entity.Policy, len(policies))
	var err error
	for i, policy := range policies {
		entities[i], err = ToPolicyEntity(policy)
		if err != nil {
			return nil, err
		}
	}
	return entities, nil
}
