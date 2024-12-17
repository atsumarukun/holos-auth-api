package entity

import (
	"crypto/rand"
	"encoding/base64"
	"holos-auth-api/internal/app/api/pkg/apierr"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var ErrUserTokenTooLong = apierr.NewApiError(http.StatusInternalServerError, "user token must be 32 characters or less")

type UserToken struct {
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
}

func NewUserToken(userID uuid.UUID) (*UserToken, error) {
	token, err := generateToken()
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

func generateToken() (string, error) {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return "", apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	token := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(buf)
	if 32 < len(token) {
		return "", ErrUserTokenTooLong
	}
	return token, nil
}
