package usecase

import (
	"context"
	"errors"
	"holos-auth-api/internal/app/api/domain"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/repository"
	"holos-auth-api/internal/app/api/domain/service"
	"holos-auth-api/internal/app/api/usecase/dto"
)

type UserUsecase interface {
	Create(context.Context, string, string, string) (*dto.UserDTO, error)
	Update(context.Context, string, string, string, string) (*dto.UserDTO, error)
	Delete(context.Context, string, string) error
}

type userUsecase struct {
	transactionObject domain.TransactionObject
	userRepository    repository.UserRepository
	userService       service.UserService
}

func NewUserUsecase(transactionObject domain.TransactionObject, userRepository repository.UserRepository, userService service.UserService) UserUsecase {
	return &userUsecase{
		transactionObject: transactionObject,
		userRepository:    userRepository,
		userService:       userService,
	}
}

func (uu *userUsecase) Create(ctx context.Context, name string, password string, confirmPassword string) (*dto.UserDTO, error) {
	user, err := entity.NewUser(name, password, confirmPassword)
	if err != nil {
		return nil, err
	}

	if err := uu.transactionObject.Transaction(ctx, func(ctx context.Context) error {
		if exists, err := uu.userService.Exists(ctx, user); err != nil {
			return err
		} else if exists {
			return errors.New("user already exists")
		}

		return uu.userRepository.Create(ctx, user)
	}); err != nil {
		return nil, err
	}

	return dto.NewUserDTO(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt), nil
}

func (uu *userUsecase) Update(ctx context.Context, name string, currentPassword string, newPassword string, confirmNewPassword string) (*dto.UserDTO, error) {
	var user *entity.User

	if err := uu.transactionObject.Transaction(ctx, func(ctx context.Context) error {
		var err error
		user, err = uu.userRepository.FindOneByName(ctx, name)
		if err != nil {
			return err
		}
		if user == nil {
			return errors.New("user not found")
		}

		if err := user.ComparePassword(currentPassword); err != nil {
			return err
		}

		if err := user.SetPassword(newPassword, confirmNewPassword); err != nil {
			return err
		}

		return uu.userRepository.Update(ctx, user)
	}); err != nil {
		return nil, err
	}

	return dto.NewUserDTO(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt), nil
}

func (uu *userUsecase) Delete(ctx context.Context, name string, password string) error {
	return uu.transactionObject.Transaction(ctx, func(ctx context.Context) error {
		user, err := uu.userRepository.FindOneByName(ctx, name)
		if err != nil {
			return err
		}
		if user == nil {
			return errors.New("user not found")
		}

		if err := user.ComparePassword(password); err != nil {
			return err
		}

		return uu.userRepository.Delete(ctx, user)
	})
}
