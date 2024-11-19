package entity

import (
	"errors"
	"holos-auth-api/internal/app/api/pkg/apierr"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNameTooShort         = apierr.NewApiError(http.StatusBadRequest, "user name must be 3 characters or more")
	ErrUserNameTooLong          = apierr.NewApiError(http.StatusBadRequest, "user name must be 24 characters or less")
	ErrInvalidUserName          = apierr.NewApiError(http.StatusBadRequest, "invalid user name")
	ErrUserPasswordDoesNotMatch = apierr.NewApiError(http.StatusBadRequest, "password does not match")
	ErrUserPasswordTooShort     = apierr.NewApiError(http.StatusBadRequest, "user password must be 8 characters or more")
	ErrUserPasswordTooLong      = apierr.NewApiError(http.StatusBadRequest, "user password must be 72 characters or less")
	ErrInvalidUserPassword      = apierr.NewApiError(http.StatusBadRequest, "invalid user password")
	ErrAuthenticationFailed     = apierr.NewApiError(http.StatusUnauthorized, "authentication failed")
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewUser(name string, password string, confirmPassword string) (*User, apierr.ApiError) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, apierr.NewApiError(http.StatusInternalServerError, err.Error())
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

func (u *User) SetName(name string) apierr.ApiError {
	if len(name) < 3 {
		return ErrUserNameTooShort
	}
	if 24 < len(name) {
		return ErrUserNameTooLong
	}
	matched, err := regexp.MatchString(`^[A-Za-z0-9_]*$`, name)
	if err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	if !matched {
		return ErrInvalidUserName
	}
	u.Name = name
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) SetPassword(password string, confirmPassword string) apierr.ApiError {
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
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	if !matched {
		return ErrInvalidUserPassword
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	u.Password = string(hashed)
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) ComparePassword(password string) apierr.ApiError {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrAuthenticationFailed
		} else {
			return apierr.NewApiError(http.StatusInternalServerError, err.Error())
		}
	}
	return nil
}
