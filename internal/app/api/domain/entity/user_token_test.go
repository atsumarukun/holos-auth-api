package entity_test

import (
	"errors"
	"holos-auth-api/internal/app/api/domain/entity"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewUserToken(t *testing.T) {
	tests := []struct {
		name        string
		inputUserID uuid.UUID
		expectError error
	}{
		{
			name:        "success",
			inputUserID: uuid.New(),
			expectError: nil,
		},
	}
	for _, tt := range tests {
		generateTime := time.Now()
		userToken, err := entity.NewUserToken(tt.inputUserID)
		if !errors.Is(err, tt.expectError) {
			t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
		}

		if tt.expectError == nil {
			if userToken.UserID != tt.inputUserID {
				t.Errorf("agent_id: expect %s but got %s", tt.inputUserID, userToken.UserID)
			}
			if len(userToken.Token) != 32 {
				t.Error("token: must be 32 characters")
			}
			if userToken.ExpiresAt.IsZero() {
				t.Error("expires_at: expect time but got empty")
			}
			if userToken.ExpiresAt.Before(generateTime.Add(time.Hour * 24 * 30)) {
				t.Error("expires_at: expect a month later")
			}
		}
	}
}
