package model

import (
	"time"

	"github.com/google/uuid"
)

type AgentTokenModel struct {
	AgentID     uuid.UUID `db:"agent_id"`
	Token       string    `db:"token"`
	GeneratedAt time.Time `db:"generated_at"`
}
