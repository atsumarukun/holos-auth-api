package entity

import (
	"crypto/rand"
	"encoding/base64"
	"holos-auth-api/internal/pkg/apierr"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var ErrUserTokenTooLong = apierr.NewApiError(http.StatusInternalServerError, "user token must be 32 characters or less")

type UserToken struct {
	UserID    uuid.UUID `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
}

func NewUserToken(userID uuid.UUID) (*UserToken, apierr.ApiError) {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return nil, apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	token := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(buf)
	if 32 < len(token) {
		return nil, ErrUserTokenTooLong
	}

	return &UserToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}, nil
}
