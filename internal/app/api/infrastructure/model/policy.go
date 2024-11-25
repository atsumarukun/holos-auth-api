package model

import (
	"time"

	"github.com/google/uuid"
)

type PolicyModel struct {
	ID             uuid.UUID `db:"id"`
	UserID         uuid.UUID `db:"user_id"`
	Name           string    `db:"name"`
	Service        string    `db:"service"`
	Path           string    `db:"path"`
	AllowedMethods string    `db:"allowed_methods"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

func NewPolicyModel(id uuid.UUID, userID uuid.UUID, name string, service string, path string, allowedMethods string, createdAt time.Time, updatedAt time.Time) *PolicyModel {
	return &PolicyModel{
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
