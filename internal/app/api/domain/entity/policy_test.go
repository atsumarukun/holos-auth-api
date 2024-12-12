package entity_test

import (
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/pkg/apierr"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestNewPolicy(t *testing.T) {
	policy, err := entity.NewPolicy(uuid.New(), "name", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}
	if policy.ID == uuid.Nil {
		t.Error("id: expect uuid but got empty")
	}
	if policy.UserID == uuid.Nil {
		t.Error("user_id: expect uuid but got empty")
	}
	if policy.Name == "" {
		t.Error("name: expect string but got empty")
	}
	if policy.Service == "" {
		t.Error("service: expect string but got empty")
	}
	if policy.Path == "" {
		t.Error("path: expect string but got empty")
	}
	if len(policy.Methods) == 0 {
		t.Error("methods: expect array but got empty")
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

func TestPolicy_SetName(t *testing.T) {
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
			expect: entity.ErrInvalidPolicyName,
		},
		{
			name:   "なまえ",
			expect: entity.ErrInvalidPolicyName,
		},
		{
			name:   strings.Repeat("a", 2),
			expect: entity.ErrPolicyNameTooShort,
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
			expect: entity.ErrPolicyNameTooLong,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy, err := entity.NewPolicy(uuid.New(), "name", "STORAGE", "/", []string{"GET"})
			if err != nil {
				t.Error(err.Error())
			}
			if err := policy.SetName(tt.name); err != tt.expect {
				if err == nil {
					t.Error("expect err but got nil")
				} else {
					t.Error(err.Error())
				}
			}
		})
	}
}

func TestPolicy_SetService(t *testing.T) {
	tests := []struct {
		name    string
		service string
		expect  apierr.ApiError
	}{
		{
			name:    "storage",
			service: "STORAGE",
			expect:  nil,
		},
		{
			name:    "content",
			service: "CONTENT",
			expect:  nil,
		},
		{
			name:    "lower_case",
			service: "storage",
			expect:  entity.ErrInvalidPolicyService,
		},
		{
			name:    "invalid",
			service: "INVALID",
			expect:  entity.ErrInvalidPolicyService,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy, err := entity.NewPolicy(uuid.New(), tt.name, "STORAGE", "/", []string{"GET"})
			if err != nil {
				t.Error(err.Error())
			}
			if err := policy.SetService(tt.service); err != tt.expect {
				if err == nil {
					t.Error("expect err but got nil")
				} else {
					t.Error(err.Error())
				}
			}
		})
	}
}

func TestPolicy_SetPath(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		expect apierr.ApiError
	}{
		{
			name:   "valid",
			path:   "/path",
			expect: nil,
		},
		{
			name:   "invalid",
			path:   "path",
			expect: entity.ErrInvalidPolicyPath,
		},
		{
			name:   "last_character_is_slash",
			path:   "/path/",
			expect: nil,
		},
		{
			name:   "max_length",
			path:   "/" + strings.Repeat("a", 254),
			expect: nil,
		},
		{
			name:   "too_long",
			path:   "/" + strings.Repeat("a", 255),
			expect: entity.ErrPolicyPathTooLong,
		},
		{
			name:   "required",
			path:   "",
			expect: entity.ErrRequiredPolicyPath,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy, err := entity.NewPolicy(uuid.New(), tt.name, "STORAGE", "/", []string{"GET"})
			if err != nil {
				t.Error(err.Error())
			}
			if err := policy.SetPath(tt.path); err != tt.expect {
				if err == nil {
					t.Error("expect err but got nil")
				} else {
					t.Error(err.Error())
				}
			}
		})
	}
}

func TestPolicy_SetMethods(t *testing.T) {
	tests := []struct {
		name    string
		Methods []string
		expect  apierr.ApiError
	}{
		{
			name:    "valid",
			Methods: []string{"GET", "POST"},
			expect:  nil,
		},
		{
			name:    "required",
			Methods: []string{},
			expect:  entity.ErrRequiredPolicyMethods,
		},
		{
			name:    "lower_case",
			Methods: []string{"get", "post"},
			expect:  entity.ErrInvalidPolicyMethods,
		},
		{
			name:    "invalid",
			Methods: []string{"INVALID"},
			expect:  entity.ErrInvalidPolicyMethods,
		},
		{
			name:    "duplication",
			Methods: []string{"GET", "POST", "GET"},
			expect:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy, err := entity.NewPolicy(uuid.New(), tt.name, "STORAGE", "/", []string{"GET"})
			if err != nil {
				t.Error(err.Error())
			}
			if err := policy.SetMethods(tt.Methods); err != tt.expect {
				if err == nil {
					t.Error("expect err but got nil")
				} else {
					t.Error(err.Error())
				}
			}
		})
	}
}

func TestPolicy_SetPermissions(t *testing.T) {
	tests := []struct {
		name   string
		effect string
		expect apierr.ApiError
	}{
		{
			name:   "allow",
			expect: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy, err := entity.NewPolicy(uuid.New(), "name", "STORAGE", "/", []string{"GET"})
			if err != nil {
				t.Error(err.Error())
			}
			agent, err := entity.NewAgent(uuid.New(), "name")
			if err != nil {
				t.Error(err.Error())
			}
			policy.SetPermissions([]*entity.Agent{agent})
		})
	}
}
