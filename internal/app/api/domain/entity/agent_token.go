package entity

import (
	"holos-auth-api/internal/app/api/domain/pkg/token"

	"github.com/google/uuid"
)

type AgentToken struct {
	AgentID uuid.UUID
	Token   string
}

func NewAgentToken(agentID uuid.UUID) (*AgentToken, error) {
	token, err := token.Generate()
	if err != nil {
		return nil, err
	}

	return &AgentToken{
		AgentID: agentID,
		Token:   token,
	}, nil
}

func RestoreAgentToken(agentID uuid.UUID, token string) *AgentToken {
	return &AgentToken{
		AgentID: agentID,
		Token:   token,
	}
}
