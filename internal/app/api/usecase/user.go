//go:generate mockgen -source=$GOFILE -destination=../../../../test/mock/usecase/$GOFILE
package usecase

import (
	"context"
	"holos-auth-api/internal/app/api/domain"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/repository"
	"holos-auth-api/internal/app/api/domain/service"
	"holos-auth-api/internal/app/api/pkg/apierr"
	"holos-auth-api/internal/app/api/usecase/dto"
	"net/http"

	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExists = apierr.NewApiError(http.StatusBadRequest, "user already exists")
	ErrUserNotFound      = apierr.NewApiError(http.StatusNotFound, "user not found")
)

type UserUsecase interface {
	Create(context.Context, string, string, string) (*dto.UserDTO, error)
	UpdateName(context.Context, uuid.UUID, string) (*dto.UserDTO, error)
	UpdatePassword(context.Context, uuid.UUID, string, string, string) (*dto.UserDTO, error)
	Delete(context.Context, uuid.UUID, string) error
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

func (u *userUsecase) Create(ctx context.Context, name string, password string, confirmPassword string) (*dto.UserDTO, error) {
	user, err := entity.NewUser(name, password, confirmPassword)
	if err != nil {
		return nil, err
	}

	if err := u.transactionObject.Transaction(ctx, func(ctx context.Context) error {
		if exists, err := u.userService.Exists(ctx, user); err != nil {
			return err
		} else if exists {
			return ErrUserAlreadyExists
		}

		return u.userRepository.Create(ctx, user)
	}); err != nil {
		return nil, err
	}

	return u.convertToDTO(user), nil
}

func (u *userUsecase) UpdateName(ctx context.Context, id uuid.UUID, name string) (*dto.UserDTO, error) {
	var user *entity.User

	if err := u.transactionObject.Transaction(ctx, func(ctx context.Context) error {
		var err error
		user, err = u.userRepository.FindOneByIDAndNotDeleted(ctx, id)
		if err != nil {
			return err
		}
		if user == nil {
			return ErrUserNotFound
		}

		if err := user.SetName(name); err != nil {
			return err
		}

		if exists, err := u.userService.Exists(ctx, user); err != nil {
			return err
		} else if exists {
			return ErrUserAlreadyExists
		}

		return u.userRepository.Update(ctx, user)
	}); err != nil {
		return nil, err
	}

	return u.convertToDTO(user), nil
}

func (u *userUsecase) UpdatePassword(ctx context.Context, id uuid.UUID, currentPassword string, newPassword string, confirmNewPassword string) (*dto.UserDTO, error) {
	var user *entity.User

	if err := u.transactionObject.Transaction(ctx, func(ctx context.Context) error {
		var err error
		user, err = u.userRepository.FindOneByIDAndNotDeleted(ctx, id)
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

		return u.userRepository.Update(ctx, user)
	}); err != nil {
		return nil, err
	}

	return u.convertToDTO(user), nil
}

func (u *userUsecase) Delete(ctx context.Context, id uuid.UUID, password string) error {
	return u.transactionObject.Transaction(ctx, func(ctx context.Context) error {
		user, err := u.userRepository.FindOneByIDAndNotDeleted(ctx, id)
		if err != nil {
			return err
		}
		if user == nil {
			return ErrUserNotFound
		}

		if err := user.ComparePassword(password); err != nil {
			return err
		}

		return u.userRepository.Delete(ctx, user)
	})
}

func (u *userUsecase) convertToDTO(user *entity.User) *dto.UserDTO {
	return dto.NewUserDTO(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt)
}
