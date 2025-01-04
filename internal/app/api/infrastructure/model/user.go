package model

import (
	"time"

	"github.com/google/uuid"
)

type UserModel struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
