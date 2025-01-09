package model

import (
	"time"

	"github.com/google/uuid"
)

type PolicyModel struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Name      string    `db:"name"`
	Effect    string    `db:"effect"`
	Service   string    `db:"service"`
	Path      string    `db:"path"`
	Methods   []byte    `db:"methods"`
	Agents    *string   `db:"agents"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
