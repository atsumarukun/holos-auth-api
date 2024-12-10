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
		effect string
		expect apierr.ApiError
	}{
		{
			name:   "allow",
			effect: "ALLOW",
			expect: nil,
		},
		{
			name:   "deny",
			effect: "DENY",
			expect: nil,
		},
		{
			name:   "invalid",
			effect: "INVALID",
			expect: entity.ErrInvalidPermissionEffect,
		},
	}
	for _, tt := range tests {
		permission, err := entity.NewPermission(uuid.New(), uuid.New(), tt.effect)
		if err != tt.expect {
			if err == nil {
				t.Error("expect err but got nil")
			} else {
				t.Error(err.Error())
			}
		}
		if tt.expect == nil {
			if permission.AgentID == uuid.Nil {
				t.Error("agent_id: expect uuid but got empty")
			}
			if permission.PolicyID == uuid.Nil {
				t.Error("policy_id: expect uuid but got empty")
			}
			if permission.Effect == "" {
				t.Error("effect: expect string but got empty")
			}
			if permission.CreatedAt.IsZero() {
				t.Error("created_at: expect time but got empty")
			}
			if permission.UpdatedAt.IsZero() {
				t.Error("updated_at: expect time but got empty")
			}
			if !permission.CreatedAt.Equal(permission.UpdatedAt) {
				t.Error("expect created_at and updated_at to be equal")
			}
		}
	}
}
