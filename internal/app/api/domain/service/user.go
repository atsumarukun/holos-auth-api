//go:generate mockgen -source=$GOFILE -destination=../../../../../test/mock/domain/service/$GOFILE
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

func (s *userService) Exists(ctx context.Context, user *entity.User) (bool, error) {
	user, err := s.userRepository.FindOneByName(ctx, user.Name)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}
