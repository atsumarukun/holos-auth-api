package model

import (
	"time"

	"github.com/google/uuid"
)

type UserTokenModel struct {
	UserID    uuid.UUID `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
}
