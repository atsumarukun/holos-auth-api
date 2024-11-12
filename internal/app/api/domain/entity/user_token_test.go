package entity_test

import (
	"holos-auth-api/internal/app/api/domain/entity"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewUserToken(t *testing.T) {
	now := time.Now()
	userToken, err := entity.NewUserToken(uuid.New())
	if err != nil {
		t.Error(err.Error())
	}
	if userToken.UserID == uuid.Nil {
		t.Error("user_id: expect uuid but got empty")
	}
	if userToken.Token == "" {
		t.Error("token: expect string but got empty")
	}
	if userToken.ExpiresAt.IsZero() {
		t.Error("expires_at: expect time but got empty")
	}
	if userToken.ExpiresAt.Before(now.Add(time.Hour * 24 * 30)) {
		t.Error("Expect expires_at a month later")
	}
}
