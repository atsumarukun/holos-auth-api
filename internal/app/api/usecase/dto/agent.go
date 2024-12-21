package dto

import (
	"time"

	"github.com/google/uuid"
)

type AgentDTO struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Policies  []uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAgentDTO(id uuid.UUID, userID uuid.UUID, name string, createdAt time.Time, updatedAt time.Time) *AgentDTO {
	return &AgentDTO{
		ID:        id,
		UserID:    userID,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
