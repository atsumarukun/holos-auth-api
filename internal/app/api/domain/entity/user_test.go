package entity_test

import (
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/pkg/apierr"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
	user, err := entity.NewUser("name", "password", "password")
	if err != nil {
		t.Error(err.Error())
	}
	if user.ID == uuid.Nil {
		t.Error("id: expect uuid but got empty")
	}
	if user.Name == "" {
		t.Error("name: expect string but got empty")
	}
	if len(user.Password) != 60 {
		t.Errorf("password: expect 60 characters but got %d characters", len(user.Password))
	}
	if user.CreatedAt.IsZero() {
		t.Error("created_at: expect time but got empty")
	}
	if user.UpdatedAt.IsZero() {
		t.Error("updated_at: expect time but got empty")
	}
	if !user.CreatedAt.Equal(user.UpdatedAt) {
		t.Error("expect created_at and updated_at to be equal")
	}
}

func TestUser_SetName(t *testing.T) {
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
			expect: entity.ErrInvalidUserName,
		},
		{
			name:   "なまえ",
			expect: entity.ErrInvalidUserName,
		},
		{
			name:   strings.Repeat("a", 2),
			expect: entity.ErrUserNameTooShort,
		},
		{
			name:   strings.Repeat("a", 3),
			expect: nil,
		},
		{
			name:   strings.Repeat("a", 24),
			expect: nil,
		},
		{
			name:   strings.Repeat("a", 25),
			expect: entity.ErrUserNameTooLong,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := entity.NewUser("name", "password", "password")
			if err != nil {
				t.Error(err.Error())
			}
			if err := user.SetName(tt.name); err != tt.expect {
				if err == nil {
					t.Error("expect err but got nil")
				} else {
					t.Error(err.Error())
				}
			}
		})
	}
}

func TestUser_SetPassword(t *testing.T) {
	tests := []struct {
		password        string
		confirmPassword string
		expect          apierr.ApiError
	}{
		{
			password:        "password",
			confirmPassword: "password",
			expect:          nil,
		},
		{
			password:        "password",
			confirmPassword: "confirm_password",
			expect:          entity.ErrUserPasswordDoesNotMatch,
		},
		{
			password:        "ぱすわーど",
			confirmPassword: "ぱすわーど",
			expect:          entity.ErrInvalidUserPassword,
		},
		{
			password:        strings.Repeat("a", 7),
			confirmPassword: strings.Repeat("a", 7),
			expect:          entity.ErrUserPasswordTooShort,
		},
		{
			password:        strings.Repeat("a", 8),
			confirmPassword: strings.Repeat("a", 8),
			expect:          nil,
		},
		{
			password:        strings.Repeat("a", 72),
			confirmPassword: strings.Repeat("a", 72),
			expect:          nil,
		},
		{
			password:        strings.Repeat("a", 73),
			confirmPassword: strings.Repeat("a", 73),
			expect:          entity.ErrUserPasswordTooLong,
		},
	}
	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			user, err := entity.NewUser("name", "password", "password")
			if err != nil {
				t.Error(err.Error())
			}
			if err := user.SetPassword(tt.password, tt.confirmPassword); err != tt.expect {
				if err == nil {
					t.Error("expect err but got nil")
				} else {
					t.Error(err.Error())
				}
			}
		})
	}
}

func TestUser_ComparePassword(t *testing.T) {
	tests := []struct {
		password string
		expect   apierr.ApiError
	}{
		{
			password: "password",
			expect:   nil,
		},
		{
			password: "test_password",
			expect:   entity.ErrAuthenticationFailed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			user, err := entity.NewUser("name", "password", "password")
			if err != nil {
				t.Error(err.Error())
			}
			if err := user.ComparePassword(tt.password); err != tt.expect {
				if err == nil {
					t.Error("expect err but got nil")
				} else {
					t.Error(err.Error())
				}
			}
		})
	}
}
