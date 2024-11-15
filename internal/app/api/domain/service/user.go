//go:generate mockgen -source=$GOFILE -destination=../../../../../test/mock/domain/service/$GOFILE
package service

import (
	"context"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/repository"
	"holos-auth-api/internal/pkg/apierr"
)

type UserService interface {
	Exists(context.Context, *entity.User) (bool, apierr.ApiError)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (us *userService) Exists(ctx context.Context, user *entity.User) (bool, apierr.ApiError) {
	user, err := us.userRepository.FindOneByName(ctx, user.Name)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}
