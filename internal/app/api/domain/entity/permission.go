package entity

import (
	"github.com/google/uuid"
)

type Permission struct {
	AgentID  uuid.UUID
	PolicyID uuid.UUID
}

func NewPermission(agentID uuid.UUID, policyID uuid.UUID) *Permission {
	return &Permission{
		AgentID:  agentID,
		PolicyID: policyID,
	}
}
