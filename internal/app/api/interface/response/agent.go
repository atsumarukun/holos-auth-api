package response

import (
	"time"

	"github.com/google/uuid"
)

type AgentResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AgentTokenResponse struct {
	GeneratedAt time.Time `json:"generated_at"`
}
