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
