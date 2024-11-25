package dto

import (
	"time"

	"github.com/google/uuid"
)

type PolicyDTO struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	Name           string
	Service        string
	Path           string
	AllowedMethods []string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewPolicyDTO(id uuid.UUID, userID uuid.UUID, name string, service string, path string, allowedMethods []string, createdAt time.Time, updatedAt time.Time) *PolicyDTO {
	return &PolicyDTO{
		ID:             id,
		UserID:         userID,
		Name:           name,
		Service:        service,
		Path:           path,
		AllowedMethods: allowedMethods,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}
