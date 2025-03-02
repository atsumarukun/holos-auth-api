package entity

import (
	"holos-auth-api/internal/app/api/domain/pkg/token"
	"time"

	"github.com/google/uuid"
)

type AgentToken struct {
	AgentID     uuid.UUID
	Token       string
	GeneratedAt time.Time
}

func NewAgentToken(agentID uuid.UUID) (*AgentToken, error) {
	token, err := token.Generate()
	if err != nil {
		return nil, err
	}

	return &AgentToken{
		AgentID:     agentID,
		Token:       token,
		GeneratedAt: time.Now(),
	}, nil
}

func RestoreAgentToken(agentID uuid.UUID, token string, generatedAt time.Time) *AgentToken {
	return &AgentToken{
		AgentID:     agentID,
		Token:       token,
		GeneratedAt: generatedAt,
	}
}
