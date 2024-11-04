package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewUser(name string, password string) (*User, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	user := &User{
		ID: id,
	}

	if err := user.SetName(name); err != nil {
		return nil, err
	}

	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	return user, nil
}

func (u *User) SetName(name string) error {
	if 255 < len(name) {
		return errors.New("id must be less than 255 characters")
	}
	u.Name = name
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) SetPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
