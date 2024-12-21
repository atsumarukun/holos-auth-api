package dto

import (
	"time"

	"github.com/google/uuid"
)

type PolicyDTO struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Service   string
	Path      string
	Methods   []string
	Agents    []uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPolicyDTO(id uuid.UUID, userID uuid.UUID, name string, service string, path string, methods []string, createdAt time.Time, updatedAt time.Time) *PolicyDTO {
	return &PolicyDTO{
		ID:        id,
		UserID:    userID,
		Name:      name,
		Service:   service,
		Path:      path,
		Methods:   methods,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
