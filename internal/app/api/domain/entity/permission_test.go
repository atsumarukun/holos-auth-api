package entity_test

import (
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/pkg/apierr"
	"testing"

	"github.com/google/uuid"
)

func TestNewPermission(t *testing.T) {
	tests := []struct {
		name   string
		expect apierr.ApiError
	}{
		{
			name:   "allow",
			expect: nil,
		},
	}
	for _, tt := range tests {
		permission := entity.NewPermission(uuid.New(), uuid.New())
		if tt.expect == nil {
			if permission.AgentID == uuid.Nil {
				t.Error("agent_id: expect uuid but got empty")
			}
			if permission.PolicyID == uuid.Nil {
				t.Error("policy_id: expect uuid but got empty")
			}
		}
	}
}
