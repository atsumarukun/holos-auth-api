package repository

import (
	"context"
	"holos-auth-api/internal/app/api/domain/entity"
)

type UserRepository interface {
	Create(context.Context, *entity.User) error
	Update(context.Context, *entity.User) error
	Delete(context.Context, *entity.User) error
	FindOneByName(context.Context, string) (*entity.User, error)
}
