package token

import (
	"crypto/rand"
	"encoding/base64"
	"holos-auth-api/internal/app/api/pkg/status"
	"net/http"
)

var ErrTokenTooLong = status.Error(http.StatusInternalServerError, "token must be 32 characters or less")

func Generate() (string, error) {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	token := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(buf)
	if 32 < len(token) {
		return "", ErrTokenTooLong
	}
	return token, nil
}
