package builder

import (
	"holos-auth-api/internal/app/api/interface/response"
	"holos-auth-api/internal/app/api/usecase/dto"
)

func ToAgentResponse(agent *dto.AgentDTO) *response.AgentResponse {
	return &response.AgentResponse{
		ID:        agent.ID,
		Name:      agent.Name,
		CreatedAt: agent.CreatedAt,
		UpdatedAt: agent.UpdatedAt,
	}
}

func ToAgentResponses(agents []*dto.AgentDTO) []*response.AgentResponse {
	responses := make([]*response.AgentResponse, len(agents))
	for i, agent := range agents {
		responses[i] = ToAgentResponse(agent)
	}
	return responses
}
