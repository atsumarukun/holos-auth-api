//go:generate mockgen -source=$GOFILE -destination=../../../../../test/mock/domain/repository/$GOFILE
package repository

import (
	"context"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/pkg/apierr"
)

type UserRepository interface {
	Create(context.Context, *entity.User) apierr.ApiError
	Update(context.Context, *entity.User) apierr.ApiError
	Delete(context.Context, *entity.User) apierr.ApiError
	FindOneByName(context.Context, string) (*entity.User, apierr.ApiError)
}
