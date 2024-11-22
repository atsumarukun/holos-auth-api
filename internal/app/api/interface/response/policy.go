package response

import (
	"time"

	"github.com/google/uuid"
)

type PolicyResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Service        string    `json:"service"`
	Path           string    `json:"path"`
	AllowedMethods []string  `json:"allowed_methods"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func NewPolicyResponse(id uuid.UUID, name string, service string, path string, allowedMethods []string, createdAt time.Time, updatedAt time.Time) *PolicyResponse {
	return &PolicyResponse{
		ID:             id,
		Name:           name,
		Service:        service,
		Path:           path,
		AllowedMethods: allowedMethods,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}
