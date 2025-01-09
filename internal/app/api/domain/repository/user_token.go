//go:generate mockgen -source=$GOFILE -destination=../../../../../test/mock/domain/repository/$GOFILE
package repository

import (
	"context"
	"holos-auth-api/internal/app/api/domain/entity"
)

type UserTokenRepository interface {
	Save(context.Context, *entity.UserToken) error
	Delete(context.Context, *entity.UserToken) error
	FindOneByTokenAndNotExpired(context.Context, string) (*entity.UserToken, error)
}
