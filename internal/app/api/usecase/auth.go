//go:generate mockgen -source=$GOFILE -destination=../../../../test/mock/usecase/$GOFILE
package usecase

import (
	"context"
	"holos-auth-api/internal/app/api/domain"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/repository"
	"holos-auth-api/internal/app/api/pkg/apierr"
	"net/http"

	"github.com/google/uuid"
)

var ErrAuthenticationFailed = apierr.NewApiError(http.StatusUnauthorized, "authentication failed")

type AuthUsecase interface {
	Signin(context.Context, string, string) (string, apierr.ApiError)
	Signout(context.Context, string) apierr.ApiError
	GetUserID(context.Context, string) (uuid.UUID, apierr.ApiError)
}

type authUsecase struct {
	transactionObject   domain.TransactionObject
	userRepository      repository.UserRepository
	userTokenRepository repository.UserTokenRepository
}

func NewAuthUsecase(transactionObject domain.TransactionObject, userRepository repository.UserRepository, userTokenRepository repository.UserTokenRepository) AuthUsecase {
	return &authUsecase{
		transactionObject:   transactionObject,
		userRepository:      userRepository,
		userTokenRepository: userTokenRepository,
	}
}

func (u *authUsecase) Signin(ctx context.Context, userName string, password string) (string, apierr.ApiError) {
	var userToken *entity.UserToken

	if err := u.transactionObject.Transaction(ctx, func(ctx context.Context) apierr.ApiError {
		user, err := u.userRepository.FindOneByName(ctx, userName)
		if err != nil {
			return err
		}
		if user == nil {
			return ErrAuthenticationFailed
		}

		if err := user.ComparePassword(password); err != nil {
			return err
		}

		userToken, err = entity.NewUserToken(user.ID)
		if err != nil {
			return err
		}

		return u.userTokenRepository.Save(ctx, userToken)
	}); err != nil {
		return "", err
	}
	return userToken.Token, nil
}

func (u *authUsecase) Signout(ctx context.Context, token string) apierr.ApiError {
	return u.transactionObject.Transaction(ctx, func(ctx context.Context) apierr.ApiError {
		userToken, err := u.userTokenRepository.FindOneByTokenAndNotExpired(ctx, token)
		if err != nil {
			return err
		}
		if userToken == nil {
			return ErrAuthenticationFailed
		}

		return u.userTokenRepository.Delete(ctx, userToken)
	})
}

func (u *authUsecase) GetUserID(ctx context.Context, token string) (uuid.UUID, apierr.ApiError) {
	userToken, err := u.userTokenRepository.FindOneByTokenAndNotExpired(ctx, token)
	if err != nil {
		return uuid.Nil, err
	}
	if userToken == nil {
		return uuid.Nil, ErrAuthenticationFailed
	}

	return userToken.UserID, nil
}
