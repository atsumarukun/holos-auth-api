//go:generate mockgen -source=$GOFILE -destination=../../../../../test/mock/domain/repository/$GOFILE
package repository

import (
	"context"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/pkg/apierr"
)

type UserTokenRepository interface {
	Save(context.Context, *entity.UserToken) apierr.ApiError
	Delete(context.Context, *entity.UserToken) apierr.ApiError
	FindOneByTokenAndNotExpired(context.Context, string) (*entity.UserToken, apierr.ApiError)
}
