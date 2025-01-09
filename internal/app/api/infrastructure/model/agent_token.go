package model

import "github.com/google/uuid"

type AgentTokenModel struct {
	AgentID uuid.UUID `db:"agent_id"`
	Token   string    `db:"token"`
}
