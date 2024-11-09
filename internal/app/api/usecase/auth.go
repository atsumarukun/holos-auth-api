//go:generate mockgen -source=$GOFILE -destination=../../../../test/mock/domain/usecase/$GOFILE
package usecase

import (
	"context"
	"holos-auth-api/internal/app/api/domain"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/repository"
	"holos-auth-api/internal/pkg/apierr"
	"net/http"
)

var ErrAuthenticationFailed = apierr.NewApiError(http.StatusUnauthorized, "authentication failed")

type AuthUsecase interface {
	Signin(context.Context, string, string) (string, apierr.ApiError)
	Signout(context.Context, string) apierr.ApiError
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

func (uu *authUsecase) Signin(ctx context.Context, userName string, password string) (string, apierr.ApiError) {
	var userToken *entity.UserToken

	if err := uu.transactionObject.Transaction(ctx, func(ctx context.Context) apierr.ApiError {
		user, err := uu.userRepository.FindOneByName(ctx, userName)
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

		return uu.userTokenRepository.Save(ctx, userToken)
	}); err != nil {
		return "", err
	}
	return userToken.Token, nil
}

func (uu *authUsecase) Signout(ctx context.Context, token string) apierr.ApiError {
	return uu.transactionObject.Transaction(ctx, func(ctx context.Context) apierr.ApiError {
		userToken, err := uu.userTokenRepository.FindOneByTokenAndNotExpired(ctx, token)
		if err != nil {
			return err
		}
		if userToken == nil {
			return ErrAuthenticationFailed
		}

		return uu.userTokenRepository.Delete(ctx, userToken)
	})
}
