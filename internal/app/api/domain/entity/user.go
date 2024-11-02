package entity

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string    `db:"id"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewUser(id string, password string) (*User, error) {
	if 255 < len(id) {
		return nil, errors.New("id must be less than 255 characters")
	}

	now := time.Now()
	user := &User{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) SetPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
