package entity_test

import (
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/pkg/apierr"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestNewAgent(t *testing.T) {
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}
	if agent.ID == uuid.Nil {
		t.Error("id: expect uuid but got empty")
	}
	if agent.UserID == uuid.Nil {
		t.Error("user_id: expect uuid but got empty")
	}
	if agent.Name == "" {
		t.Error("name: expect string but got empty")
	}
	if agent.CreatedAt.IsZero() {
		t.Error("created_at: expect time but got empty")
	}
	if agent.UpdatedAt.IsZero() {
		t.Error("updated_at: expect time but got empty")
	}
	if !agent.CreatedAt.Equal(agent.UpdatedAt) {
		t.Error("expect created_at and updated_at to be equal")
	}
}

func TestAgent_SetName(t *testing.T) {
	tests := []struct {
		name   string
		expect apierr.ApiError
	}{
		{
			name:   "valid_NAME",
			expect: nil,
		},
		{
			name:   "invalid-NAME",
			expect: entity.ErrInvalidAgentName,
		},
		{
			name:   "なまえ",
			expect: entity.ErrInvalidAgentName,
		},
		{
			name:   strings.Repeat("a", 2),
			expect: entity.ErrAgentNameTooShort,
		},
		{
			name:   strings.Repeat("a", 3),
			expect: nil,
		},
		{
			name:   strings.Repeat("a", 255),
			expect: nil,
		},
		{
			name:   strings.Repeat("a", 256),
			expect: entity.ErrAgentNameTooLong,
		},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			agent, err := entity.NewAgent(uuid.New(), "name")
			if err != nil {
				t.Error(err.Error())
			}
			if err := agent.SetName(tt.name); err != tt.expect {
				if err == nil {
					t.Error("expect err but got nil")
				} else {
					t.Error(err.Error())
				}
			}
		})
	}
}

func TestAgent_SetPermissions(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			agent, err := entity.NewAgent(uuid.New(), "name")
			if err != nil {
				t.Error(err.Error())
			}
			policy, err := entity.NewPolicy(uuid.New(), "name", "STORAGE", "/", []string{"GET"})
			if err != nil {
				t.Error(err.Error())
			}
			agent.SetPermissions([]*entity.Policy{policy})
		})
	}
}
