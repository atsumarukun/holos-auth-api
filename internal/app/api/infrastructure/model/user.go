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

func NewUserModel(id uuid.UUID, name string, password string, createdAt time.Time, updatedAt time.Time) *UserModel {
	return &UserModel{
		ID:        id,
		Name:      name,
		Password:  password,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
