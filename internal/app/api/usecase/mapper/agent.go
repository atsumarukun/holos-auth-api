package mapper

import (
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/usecase/dto"
)

func ToAgentDTO(agent *entity.Agent) *dto.AgentDTO {
	return &dto.AgentDTO{
		ID:        agent.ID,
		UserID:    agent.UserID,
		Name:      agent.Name,
		Policies:  agent.Policies,
		CreatedAt: agent.CreatedAt,
		UpdatedAt: agent.UpdatedAt,
	}
}

func ToAgentDTOs(agents []*entity.Agent) []*dto.AgentDTO {
	dtos := make([]*dto.AgentDTO, len(agents))
	for i, agent := range agents {
		dtos[i] = ToAgentDTO(agent)
	}
	return dtos
}
