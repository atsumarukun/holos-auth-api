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

type AgentTokenDTO struct {
	AgentID     uuid.UUID
	Token       string
	GeneratedAt time.Time
}
