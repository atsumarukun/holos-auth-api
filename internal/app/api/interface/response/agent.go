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

func NewAgentResponse(id uuid.UUID, name string, createdAt time.Time, updatedAt time.Time) *AgentResponse {
	return &AgentResponse{
		ID:        id,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
