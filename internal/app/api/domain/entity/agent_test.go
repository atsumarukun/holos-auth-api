package entity_test

import (
	"errors"
	"holos-auth-api/internal/app/api/domain/entity"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestNewAgent(t *testing.T) {
	tests := []struct {
		name        string
		inputUserID uuid.UUID
		inputName   string
		expectError error
	}{
		{
			name:        "success",
			inputUserID: uuid.New(),
			inputName:   "name",
			expectError: nil,
		},
		{
			name:        "invalid name",
			inputUserID: uuid.New(),
			inputName:   "なまえ",
			expectError: entity.ErrInvalidAgentName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := entity.NewAgent(tt.inputUserID, tt.inputName)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if tt.expectError == nil {
				if agent.ID == uuid.Nil {
					t.Error("id: expect uuid but got empty")
				}
				if agent.UserID != tt.inputUserID {
					t.Errorf("user_id: expect %v but got %v", tt.inputUserID, agent.UserID)
				}
				if agent.Name != tt.inputName {
					t.Errorf("name: expect %s but got %s", tt.inputName, agent.Name)
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
		})
	}
}

func TestAgent_SetName(t *testing.T) {
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name        string
		inputName   string
		expectError error
	}{
		{
			name:        "lower case",
			inputName:   "name",
			expectError: nil,
		},
		{
			name:        "upper case",
			inputName:   "NAME",
			expectError: nil,
		},
		{
			name:        "numeral",
			inputName:   "name01",
			expectError: nil,
		},
		{
			name:        "underscore",
			inputName:   "sample_name",
			expectError: nil,
		},
		{
			name:        "hiragana",
			inputName:   "なまえ",
			expectError: entity.ErrInvalidAgentName,
		},
		{
			name:        "hyphen",
			inputName:   "sample-name",
			expectError: entity.ErrInvalidAgentName,
		},
		{
			name:        "slash",
			inputName:   "sample/name",
			expectError: entity.ErrInvalidAgentName,
		},
		{
			name:        "2 characters",
			inputName:   strings.Repeat("a", 2),
			expectError: entity.ErrAgentNameTooShort,
		},
		{
			name:        "3 characters",
			inputName:   strings.Repeat("a", 3),
			expectError: nil,
		},
		{
			name:        "255 characters",
			inputName:   strings.Repeat("a", 255),
			expectError: nil,
		},
		{
			name:        "256 characters",
			inputName:   strings.Repeat("a", 256),
			expectError: entity.ErrAgentNameTooLong,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedAt := agent.UpdatedAt
			if err := agent.SetName(tt.inputName); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if tt.expectError == nil {
				if !agent.UpdatedAt.After(updatedAt) {
					t.Error("updatedAt has not been updated")
				}
			}
		})
	}
}

func TestAgent_SetPolicies(t *testing.T) {
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}
	policy, err := entity.NewPolicy(uuid.New(), "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name          string
		inputPolicies []*entity.Policy
		expectResult  []uuid.UUID
		expectError   error
	}{
		{
			name:          "success",
			inputPolicies: []*entity.Policy{policy},
			expectResult:  []uuid.UUID{policy.ID},
			expectError:   nil,
		},
		{
			name:          "empty",
			inputPolicies: []*entity.Policy{},
			expectResult:  []uuid.UUID{},
			expectError:   nil,
		},
		{
			name:          "duplication",
			inputPolicies: []*entity.Policy{policy, policy},
			expectResult:  []uuid.UUID{policy.ID},
			expectError:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedAt := agent.UpdatedAt
			agent.SetPolicies(tt.inputPolicies)
			if diff := cmp.Diff(agent.Policies, tt.expectResult); diff != "" {
				t.Error(diff)
			}
			if !agent.UpdatedAt.After(updatedAt) {
				t.Error("updatedAt has not been updated")
			}
		})
	}
}
