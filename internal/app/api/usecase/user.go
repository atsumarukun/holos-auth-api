//go:generate mockgen -source=$GOFILE -destination=../../../../test/mock/usecase/$GOFILE
package usecase

import (
	"context"
	"holos-auth-api/internal/app/api/domain"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/repository"
	"holos-auth-api/internal/app/api/domain/service"
	"holos-auth-api/internal/app/api/usecase/dto"
	"holos-auth-api/internal/pkg/apierr"
	"net/http"
)

var (
	ErrUserAlreadyExists = apierr.NewApiError(http.StatusBadRequest, "user already exists")
	ErrUserNotFound      = apierr.NewApiError(http.StatusNotFound, "user not found")
)

type UserUsecase interface {
	Create(context.Context, string, string, string) (*dto.UserDTO, apierr.ApiError)
	Update(context.Context, string, string, string, string) (*dto.UserDTO, apierr.ApiError)
	Delete(context.Context, string, string) apierr.ApiError
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

func (uu *userUsecase) Create(ctx context.Context, name string, password string, confirmPassword string) (*dto.UserDTO, apierr.ApiError) {
	user, err := entity.NewUser(name, password, confirmPassword)
	if err != nil {
		return nil, err
	}

	if err := uu.transactionObject.Transaction(ctx, func(ctx context.Context) apierr.ApiError {
		if exists, err := uu.userService.Exists(ctx, user); err != nil {
			return err
		} else if exists {
			return ErrUserAlreadyExists
		}

		return uu.userRepository.Create(ctx, user)
	}); err != nil {
		return nil, err
	}

	return dto.NewUserDTO(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt), nil
}

func (uu *userUsecase) Update(ctx context.Context, name string, currentPassword string, newPassword string, confirmNewPassword string) (*dto.UserDTO, apierr.ApiError) {
	var user *entity.User

	if err := uu.transactionObject.Transaction(ctx, func(ctx context.Context) apierr.ApiError {
		var err apierr.ApiError
		user, err = uu.userRepository.FindOneByName(ctx, name)
		if err != nil {
			return err
		}
		if user == nil {
			return ErrUserNotFound
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

func (uu *userUsecase) Delete(ctx context.Context, name string, password string) apierr.ApiError {
	return uu.transactionObject.Transaction(ctx, func(ctx context.Context) apierr.ApiError {
		user, err := uu.userRepository.FindOneByName(ctx, name)
		if err != nil {
			return err
		}
		if user == nil {
			return ErrUserNotFound
		}

		if err := user.ComparePassword(password); err != nil {
			return err
		}

		return uu.userRepository.Delete(ctx, user)
	})
}
