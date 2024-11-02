package service

import (
	"context"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/repository"
)

type UserService interface {
	Exists(context.Context, *entity.User) (bool, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (us *userService) Exists(ctx context.Context, user *entity.User) (bool, error) {
	u, err := us.userRepository.FindOneByID(ctx, user.ID)
	if err != nil {
		return false, nil
	}
	return u != nil, nil
}
