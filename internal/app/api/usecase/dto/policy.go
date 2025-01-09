package dto

import (
	"time"

	"github.com/google/uuid"
)

type PolicyDTO struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Effect    string
	Service   string
	Path      string
	Methods   []string
	Agents    []uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
