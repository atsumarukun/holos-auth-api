package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserDTO struct {
	ID        uuid.UUID
	Name      string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUserDTO(id uuid.UUID, name string, password string, createdAt time.Time, updatedAt time.Time) *UserDTO {
	return &UserDTO{
		ID:        id,
		Name:      name,
		Password:  password,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
