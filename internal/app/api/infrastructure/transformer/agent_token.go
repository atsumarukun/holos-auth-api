package transformer

import (
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/infrastructure/model"
)

func ToAgentTokenModel(agentToken *entity.AgentToken) *model.AgentTokenModel {
	return &model.AgentTokenModel{
		AgentID:     agentToken.AgentID,
		Token:       agentToken.Token,
		GeneratedAt: agentToken.GeneratedAt,
	}
}

func ToAgentTokenEntity(agentToken *model.AgentTokenModel) *entity.AgentToken {
	return entity.RestoreAgentToken(
		agentToken.AgentID,
		agentToken.Token,
		agentToken.GeneratedAt,
	)
}
