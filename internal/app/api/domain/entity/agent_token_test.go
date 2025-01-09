package entity_test

import (
	"errors"
	"holos-auth-api/internal/app/api/domain/entity"
	"testing"

	"github.com/google/uuid"
)

func TestNewAgentToken(t *testing.T) {
	tests := []struct {
		name         string
		inputAgentID uuid.UUID
		expectError  error
	}{
		{
			name:         "success",
			inputAgentID: uuid.New(),
			expectError:  nil,
		},
	}
	for _, tt := range tests {
		agentToken, err := entity.NewAgentToken(tt.inputAgentID)
		if !errors.Is(err, tt.expectError) {
			t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
		}

		if tt.expectError == nil {
			if agentToken.AgentID != tt.inputAgentID {
				t.Errorf("agent_id: expect %s but got %s", tt.inputAgentID, agentToken.AgentID)
			}
			if len(agentToken.Token) != 32 {
				t.Error("token: must be 32 characters")
			}
		}
	}
}
