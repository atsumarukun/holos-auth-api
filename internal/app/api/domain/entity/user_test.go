package entity_test

import (
	"errors"
	"holos-auth-api/internal/app/api/domain/entity"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name                 string
		inputName            string
		inputPassword        string
		inputConfirmPassword string
		expectError          error
	}{
		{
			name:                 "success",
			inputName:            "name",
			inputPassword:        "password",
			inputConfirmPassword: "password",
			expectError:          nil,
		},
		{
			name:                 "invalid name",
			inputName:            "なまえ",
			inputPassword:        "password",
			inputConfirmPassword: "password",
			expectError:          entity.ErrInvalidUserName,
		},
		{
			name:                 "invalid password",
			inputName:            "name",
			inputPassword:        "ぱすわーど",
			inputConfirmPassword: "ぱすわーど",
			expectError:          entity.ErrInvalidUserPassword,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := entity.NewUser(tt.inputName, tt.inputPassword, tt.inputConfirmPassword)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if tt.expectError == nil {
				if user.ID == uuid.Nil {
					t.Error("id: expect uuid but got empty")
				}
				if user.Name != tt.inputName {
					t.Errorf("name: expect string but got %s", tt.inputName)
				}
				if len(user.Password) != 60 {
					t.Errorf("password: expect 60 characters but got %d characters", len(user.Password))
				}
				if user.Password == tt.inputPassword {
					t.Error("password: expect hashed text but got plain text")
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
		})
	}
}

func TestUser_SetName(t *testing.T) {
	user, err := entity.NewUser("name", "password", "password")
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
			expectError: entity.ErrInvalidUserName,
		},
		{
			name:        "hyphen",
			inputName:   "sample-name",
			expectError: entity.ErrInvalidUserName,
		},
		{
			name:        "slash",
			inputName:   "sample/name",
			expectError: entity.ErrInvalidUserName,
		},
		{
			name:        "2 characters",
			inputName:   strings.Repeat("a", 2),
			expectError: entity.ErrUserNameTooShort,
		},
		{
			name:        "3 characters",
			inputName:   strings.Repeat("a", 3),
			expectError: nil,
		},
		{
			name:        "24 characters",
			inputName:   strings.Repeat("a", 24),
			expectError: nil,
		},
		{
			name:        "25 characters",
			inputName:   strings.Repeat("a", 25),
			expectError: entity.ErrUserNameTooLong,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedAt := user.UpdatedAt
			if err := user.SetName(tt.inputName); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if tt.expectError == nil {
				if !user.UpdatedAt.After(updatedAt) {
					t.Error("updatedAt has not been updated")
				}
			}
		})
	}
}

func TestUser_SetPassword(t *testing.T) {
	user, err := entity.NewUser("name", "password", "password")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name                 string
		inputPassword        string
		inputConfirmPassword string
		expectError          error
	}{
		{
			name:                 "lower case",
			inputPassword:        "password",
			inputConfirmPassword: "password",
			expectError:          nil,
		},
		{
			name:                 "upper case",
			inputPassword:        "PASSWORD",
			inputConfirmPassword: "PASSWORD",
			expectError:          nil,
		},
		{
			name:                 "numeral",
			inputPassword:        "password01",
			inputConfirmPassword: "password01",
			expectError:          nil,
		},
		{
			name:                 "symbol",
			inputPassword:        "!@#$%^&*()_-+=[]{};:'\",.<>?/|~",
			inputConfirmPassword: "!@#$%^&*()_-+=[]{};:'\",.<>?/|~",
			expectError:          nil,
		},
		{
			name:                 "hiragana",
			inputPassword:        "ぱすわーど",
			inputConfirmPassword: "ぱすわーど",
			expectError:          entity.ErrInvalidUserPassword,
		},
		{
			name:                 "7 characters",
			inputPassword:        strings.Repeat("a", 7),
			inputConfirmPassword: strings.Repeat("a", 7),
			expectError:          entity.ErrUserPasswordTooShort,
		},
		{
			name:                 "8 characters",
			inputPassword:        strings.Repeat("a", 8),
			inputConfirmPassword: strings.Repeat("a", 8),
			expectError:          nil,
		},
		{
			name:                 "72 characters",
			inputPassword:        strings.Repeat("a", 72),
			inputConfirmPassword: strings.Repeat("a", 72),
			expectError:          nil,
		},
		{
			name:                 "73 characters",
			inputPassword:        strings.Repeat("a", 73),
			inputConfirmPassword: strings.Repeat("a", 73),
			expectError:          entity.ErrUserPasswordTooLong,
		},
		{
			name:                 "password does not match",
			inputPassword:        "password",
			inputConfirmPassword: "confirm_password",
			expectError:          entity.ErrUserPasswordDoesNotMatch,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedAt := user.UpdatedAt
			if err := user.SetPassword(tt.inputPassword, tt.inputConfirmPassword); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if tt.expectError == nil {
				if !user.UpdatedAt.After(updatedAt) {
					t.Error("updatedAt has not been updated")
				}
			}
		})
	}
}

func TestUser_ComparePassword(t *testing.T) {
	user, err := entity.NewUser("name", "password", "password")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name          string
		inputPassword string
		expectError   error
	}{
		{
			name:          "success",
			inputPassword: "password",
			expectError:   nil,
		},
		{
			name:          "failure",
			inputPassword: "PASSWORD",
			expectError:   entity.ErrAuthenticationFailed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := user.ComparePassword(tt.inputPassword); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
		})
	}
}
