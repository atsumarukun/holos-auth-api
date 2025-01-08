package entity_test

import (
	"errors"
	"holos-auth-api/internal/app/api/domain/entity"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestNewPolicy(t *testing.T) {
	tests := []struct {
		name         string
		inputUserID  uuid.UUID
		inputName    string
		inputEffect  string
		inputService string
		inputPath    string
		inputMethods []string
		expectError  error
	}{
		{
			name:         "success",
			inputUserID:  uuid.New(),
			inputName:    "name",
			inputEffect:  "ALLOW",
			inputService: "STORAGE",
			inputPath:    "/",
			inputMethods: []string{"GET"},
			expectError:  nil,
		},
		{
			name:         "invalied name",
			inputUserID:  uuid.New(),
			inputName:    "なまえ",
			inputEffect:  "ALLOW",
			inputService: "STORAGE",
			inputPath:    "/",
			inputMethods: []string{"GET"},
			expectError:  entity.ErrInvalidPolicyName,
		},
		{
			name:         "invalid effect",
			inputUserID:  uuid.New(),
			inputName:    "name",
			inputEffect:  "EFFECT",
			inputService: "STORAGE",
			inputPath:    "/",
			inputMethods: []string{"GET"},
			expectError:  entity.ErrInvalidPolicyEffect,
		},
		{
			name:         "invalid service",
			inputUserID:  uuid.New(),
			inputName:    "name",
			inputEffect:  "ALLOW",
			inputService: "SERVICE",
			inputPath:    "/",
			inputMethods: []string{"GET"},
			expectError:  entity.ErrInvalidPolicyService,
		},
		{
			name:         "invalid path",
			inputUserID:  uuid.New(),
			inputName:    "name",
			inputEffect:  "ALLOW",
			inputService: "STORAGE",
			inputPath:    "path",
			inputMethods: []string{"GET"},
			expectError:  entity.ErrInvalidPolicyPath,
		},
		{
			name:         "invalid methods",
			inputUserID:  uuid.New(),
			inputName:    "name",
			inputEffect:  "ALLOW",
			inputService: "STORAGE",
			inputPath:    "/",
			inputMethods: []string{"PATCH"},
			expectError:  entity.ErrInvalidPolicyMethods,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy, err := entity.NewPolicy(tt.inputUserID, tt.inputName, tt.inputEffect, tt.inputService, tt.inputPath, tt.inputMethods)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if tt.expectError == nil {
				if policy.ID == uuid.Nil {
					t.Error("id: expect uuid but got empty")
				}
				if policy.UserID != tt.inputUserID {
					t.Errorf("user_id: expect %v but got %v", tt.inputUserID, policy.UserID)
				}
				if policy.Name != tt.inputName {
					t.Errorf("name: expect %s but got %s", tt.inputName, policy.Name)
				}
				if policy.Effect != tt.inputEffect {
					t.Errorf("effect: expect %s but got %s", tt.inputEffect, policy.Effect)
				}
				if policy.Service != tt.inputService {
					t.Errorf("service: expect %s but got %s", tt.inputService, policy.Service)
				}
				if policy.Path != tt.inputPath {
					t.Errorf("path: expect %s but got %s", tt.inputPath, policy.Path)
				}
				if diff := cmp.Diff(policy.Methods, tt.inputMethods); diff != "" {
					t.Error(diff)
				}
				if policy.CreatedAt.IsZero() {
					t.Error("created_at: expect time but got empty")
				}
				if policy.UpdatedAt.IsZero() {
					t.Error("updated_at: expect time but got empty")
				}
				if !policy.CreatedAt.Equal(policy.UpdatedAt) {
					t.Error("expect created_at and updated_at to be equal")
				}
			}
		})
	}
}

func TestPolicy_SetName(t *testing.T) {
	policy, err := entity.NewPolicy(uuid.New(), "name", "ALLOW", "STORAGE", "/", []string{"GET"})
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
			expectError: entity.ErrInvalidPolicyName,
		},
		{
			name:        "hyphen",
			inputName:   "sample-name",
			expectError: entity.ErrInvalidPolicyName,
		},
		{
			name:        "slash",
			inputName:   "sample/name",
			expectError: entity.ErrInvalidPolicyName,
		},
		{
			name:        "2 characters",
			inputName:   strings.Repeat("a", 2),
			expectError: entity.ErrPolicyNameTooShort,
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
			expectError: entity.ErrPolicyNameTooLong,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedAt := policy.UpdatedAt
			if err := policy.SetName(tt.inputName); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if tt.expectError == nil {
				if !policy.UpdatedAt.After(updatedAt) {
					t.Error("updatedAt has not been updated")
				}
			}
		})
	}
}

func TestPolicy_SetEffect(t *testing.T) {
	policy, err := entity.NewPolicy(uuid.New(), "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name        string
		inputEffect string
		expectError error
	}{
		{
			name:        "allow",
			inputEffect: "ALLOW",
			expectError: nil,
		},
		{
			name:        "deny",
			inputEffect: "DENY",
			expectError: nil,
		},
		{
			name:        "lower case",
			inputEffect: "allow",
			expectError: entity.ErrInvalidPolicyEffect,
		},
		{
			name:        "no service",
			inputEffect: "EFFECT",
			expectError: entity.ErrInvalidPolicyEffect,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedAt := policy.UpdatedAt
			if err := policy.SetEffect(tt.inputEffect); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if tt.expectError == nil {
				if !policy.UpdatedAt.After(updatedAt) {
					t.Error("updatedAt has not been updated")
				}
			}
		})
	}
}

func TestPolicy_SetService(t *testing.T) {
	policy, err := entity.NewPolicy(uuid.New(), "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name         string
		inputService string
		expectError  error
	}{
		{
			name:         "storage",
			inputService: "STORAGE",
			expectError:  nil,
		},
		{
			name:         "content",
			inputService: "CONTENT",
			expectError:  nil,
		},
		{
			name:         "lower case",
			inputService: "storage",
			expectError:  entity.ErrInvalidPolicyService,
		},
		{
			name:         "no service",
			inputService: "SERVICE",
			expectError:  entity.ErrInvalidPolicyService,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedAt := policy.UpdatedAt
			if err := policy.SetService(tt.inputService); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if tt.expectError == nil {
				if !policy.UpdatedAt.After(updatedAt) {
					t.Error("updatedAt has not been updated")
				}
			}
		})
	}
}

func TestPolicy_SetPath(t *testing.T) {
	policy, err := entity.NewPolicy(uuid.New(), "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name        string
		inputPath   string
		expectError error
	}{
		{
			name:        "lower case",
			inputPath:   "/path",
			expectError: nil,
		},
		{
			name:        "upper case",
			inputPath:   "/PATH",
			expectError: entity.ErrInvalidPolicyPath,
		},
		{
			name:        "numeral",
			inputPath:   "/path01",
			expectError: entity.ErrInvalidPolicyPath,
		},
		{
			name:        "underscore",
			inputPath:   "/sample_path",
			expectError: entity.ErrInvalidPolicyPath,
		},
		{
			name:        "hiragana",
			inputPath:   "/ぱす",
			expectError: entity.ErrInvalidPolicyPath,
		},
		{
			name:        "hyphen",
			inputPath:   "/sample-path",
			expectError: nil,
		},
		{
			name:        "colon",
			inputPath:   "/:id",
			expectError: nil,
		},
		{
			name:        "slash only",
			inputPath:   "/",
			expectError: nil,
		},
		{
			name:        "not slash start",
			inputPath:   "path",
			expectError: entity.ErrInvalidPolicyPath,
		},
		{
			name:        "slash end",
			inputPath:   "/path/",
			expectError: entity.ErrInvalidPolicyPath,
		},
		{
			name:        "255 characters",
			inputPath:   "/" + strings.Repeat("a", 254),
			expectError: nil,
		},
		{
			name:        "256 characters",
			inputPath:   "/" + strings.Repeat("a", 255),
			expectError: entity.ErrPolicyPathTooLong,
		},
		{
			name:        "empty",
			inputPath:   "",
			expectError: entity.ErrRequiredPolicyPath,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedAt := policy.UpdatedAt
			if err := policy.SetPath(tt.inputPath); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if tt.expectError == nil {
				if !policy.UpdatedAt.After(updatedAt) {
					t.Error("updatedAt has not been updated")
				}
			}
		})
	}
}

func TestPolicy_SetMethods(t *testing.T) {
	policy, err := entity.NewPolicy(uuid.New(), "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name         string
		inputMethods []string
		expectError  error
	}{
		{
			name:         "allow methods",
			inputMethods: []string{"GET", "POST", "PUT", "DELETE"},
			expectError:  nil,
		},
		{
			name:         "head",
			inputMethods: []string{"HEAD"},
			expectError:  entity.ErrInvalidPolicyMethods,
		},
		{
			name:         "connect",
			inputMethods: []string{"CONNECT"},
			expectError:  entity.ErrInvalidPolicyMethods,
		},
		{
			name:         "options",
			inputMethods: []string{"OPTIONS"},
			expectError:  entity.ErrInvalidPolicyMethods,
		},
		{
			name:         "trace",
			inputMethods: []string{"TRACE"},
			expectError:  entity.ErrInvalidPolicyMethods,
		},
		{
			name:         "patch",
			inputMethods: []string{"PATCH"},
			expectError:  entity.ErrInvalidPolicyMethods,
		},
		{
			name:         "lower case",
			inputMethods: []string{"get"},
			expectError:  entity.ErrInvalidPolicyMethods,
		},
		{
			name:         "duplication",
			inputMethods: []string{"GET", "POST", "GET"},
			expectError:  nil,
		},
		{
			name:         "empty",
			inputMethods: []string{},
			expectError:  entity.ErrRequiredPolicyMethods,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedAt := policy.UpdatedAt
			if err := policy.SetMethods(tt.inputMethods); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if tt.expectError == nil {
				if !policy.UpdatedAt.After(updatedAt) {
					t.Error("updatedAt has not been updated")
				}
			}
		})
	}
}

func TestPolicy_SetAgents(t *testing.T) {
	policy, err := entity.NewPolicy(uuid.New(), "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}
	agent, err := entity.NewAgent(policy.UserID, "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name         string
		inputAgents  []*entity.Agent
		expectResult []uuid.UUID
		expectError  error
	}{
		{
			name:         "success",
			inputAgents:  []*entity.Agent{agent},
			expectResult: []uuid.UUID{agent.ID},
			expectError:  nil,
		},
		{
			name:         "empty",
			inputAgents:  []*entity.Agent{},
			expectResult: []uuid.UUID{},
			expectError:  nil,
		},
		{
			name:         "duplication",
			inputAgents:  []*entity.Agent{agent, agent},
			expectResult: []uuid.UUID{agent.ID},
			expectError:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedAt := policy.UpdatedAt
			policy.SetAgents(tt.inputAgents)
			if diff := cmp.Diff(policy.Agents, tt.expectResult); diff != "" {
				t.Error(diff)
			}
			if !policy.UpdatedAt.After(updatedAt) {
				t.Error("updatedAt has not been updated")
			}
		})
	}
}
