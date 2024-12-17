package entity

import (
	"errors"
	"holos-auth-api/internal/app/api/pkg/status"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNameTooShort         = status.Error(http.StatusBadRequest, "user name must be 3 characters or more")
	ErrUserNameTooLong          = status.Error(http.StatusBadRequest, "user name must be 24 characters or less")
	ErrInvalidUserName          = status.Error(http.StatusBadRequest, "invalid user name")
	ErrUserPasswordDoesNotMatch = status.Error(http.StatusBadRequest, "password does not match")
	ErrUserPasswordTooShort     = status.Error(http.StatusBadRequest, "user password must be 8 characters or more")
	ErrUserPasswordTooLong      = status.Error(http.StatusBadRequest, "user password must be 72 characters or less")
	ErrInvalidUserPassword      = status.Error(http.StatusBadRequest, "invalid user password")
	ErrAuthenticationFailed     = status.Error(http.StatusUnauthorized, "authentication failed")
)

type User struct {
	ID        uuid.UUID
	Name      string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(name string, password string, confirmPassword string) (*User, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	user := &User{
		ID: id,
	}

	if err := user.SetName(name); err != nil {
		return nil, err
	}

	if err := user.SetPassword(password, confirmPassword); err != nil {
		return nil, err
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	return user, nil
}

func RestoreUser(id uuid.UUID, name string, password string, createdAt time.Time, updatedAt time.Time) *User {
	return &User{
		ID:        id,
		Name:      name,
		Password:  password,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (u *User) SetName(name string) error {
	if len(name) < 3 {
		return ErrUserNameTooShort
	}
	if 24 < len(name) {
		return ErrUserNameTooLong
	}
	matched, err := regexp.MatchString(`^[A-Za-z0-9_]*$`, name)
	if err != nil {
		return status.Error(http.StatusInternalServerError, err.Error())
	}
	if !matched {
		return ErrInvalidUserName
	}
	u.Name = name
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) SetPassword(password string, confirmPassword string) error {
	if password != confirmPassword {
		return ErrUserPasswordDoesNotMatch
	}
	if len(password) < 8 {
		return ErrUserPasswordTooShort
	}
	if 72 < len(password) {
		return ErrUserPasswordTooLong
	}
	matched, err := regexp.MatchString(`^[A-Za-z0-9!@#$%^&*()_\-+=\[\]{};:'",.<>?/\\|~]*$`, password)
	if err != nil {
		return status.Error(http.StatusInternalServerError, err.Error())
	}
	if !matched {
		return ErrInvalidUserPassword
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return status.Error(http.StatusInternalServerError, err.Error())
	}
	u.Password = string(hashed)
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) ComparePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrAuthenticationFailed
		} else {
			return status.Error(http.StatusInternalServerError, err.Error())
		}
	}
	return nil
}
