package model

import (
	"time"

	"github.com/google/uuid"
)

type AgentModel struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewAgentModel(id uuid.UUID, userID uuid.UUID, name string, createdAt time.Time, updatedAt time.Time) *AgentModel {
	return &AgentModel{
		ID:        id,
		UserID:    userID,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
