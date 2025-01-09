package entity

import (
	"holos-auth-api/internal/app/api/domain/pkg/token"
	"time"

	"github.com/google/uuid"
)

type UserToken struct {
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
}

func NewUserToken(userID uuid.UUID) (*UserToken, error) {
	token, err := token.Generate()
	if err != nil {
		return nil, err
	}

	return &UserToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}, nil
}

func RestoreUserToken(userID uuid.UUID, token string, expiresAt time.Time) *UserToken {
	return &UserToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
}
