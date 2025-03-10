//go:generate mockgen -source=$GOFILE -destination=../../../../test/mock/usecase/$GOFILE
package usecase

import (
	"context"
	"holos-auth-api/internal/app/api/domain"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/repository"
	"holos-auth-api/internal/app/api/domain/service"
	"holos-auth-api/internal/app/api/pkg/status"
	"net/http"

	"github.com/google/uuid"
)

var (
	ErrAuthenticationFailed = status.Error(http.StatusUnauthorized, "authentication failed")
	ErrAuthorizationFaild   = status.Error(http.StatusForbidden, "authorization failed")
)

type AuthUsecase interface {
	Signin(context.Context, string, string) (string, error)
	Signout(context.Context, string) error
	Authenticate(context.Context, string) (uuid.UUID, error)
	Authorize(context.Context, string, string, string, string, string) (uuid.UUID, error)
}

type authUsecase struct {
	transactionObject   domain.TransactionObject
	userRepository      repository.UserRepository
	userTokenRepository repository.UserTokenRepository
	agentRepository     repository.AgentRepository
	agentService        service.AgentService
}

func NewAuthUsecase(
	transactionObject domain.TransactionObject,
	userRepository repository.UserRepository,
	userTokenRepository repository.UserTokenRepository,
	agentRepository repository.AgentRepository,
	agentService service.AgentService,
) AuthUsecase {
	return &authUsecase{
		transactionObject:   transactionObject,
		userRepository:      userRepository,
		userTokenRepository: userTokenRepository,
		agentRepository:     agentRepository,
		agentService:        agentService,
	}
}

func (u *authUsecase) Signin(ctx context.Context, userName string, password string) (string, error) {
	var userToken *entity.UserToken

	if err := u.transactionObject.Transaction(ctx, func(ctx context.Context) error {
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

func (u *authUsecase) Signout(ctx context.Context, token string) error {
	return u.transactionObject.Transaction(ctx, func(ctx context.Context) error {
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

func (u *authUsecase) Authenticate(ctx context.Context, token string) (uuid.UUID, error) {
	userToken, err := u.userTokenRepository.FindOneByTokenAndNotExpired(ctx, token)
	if err != nil {
		return uuid.Nil, err
	}
	if userToken == nil {
		return uuid.Nil, ErrAuthenticationFailed
	}

	return userToken.UserID, nil
}

func (u *authUsecase) Authorize(ctx context.Context, token string, operatorType string, service string, path string, method string) (uuid.UUID, error) {
	switch operatorType {
	case "USER":
		return u.Authenticate(ctx, token)
	case "AGENT":
		var userID uuid.UUID
		if err := u.transactionObject.Transaction(ctx, func(ctx context.Context) error {
			agent, err := u.agentRepository.FindOneByTokenAndNotDeleted(ctx, token)
			if err != nil {
				return err
			}
			if agent == nil {
				return ErrAuthenticationFailed
			}

			hasPermission, err := u.agentService.HasPermission(ctx, agent, service, path, method)
			if err != nil {
				return err
			}
			if hasPermission {
				userID = agent.UserID
			}

			return nil
		}); err != nil {
			return uuid.Nil, err
		}

		if userID != uuid.Nil {
			return userID, nil
		}
		return uuid.Nil, ErrAuthorizationFaild
	default:
		return uuid.Nil, ErrAuthenticationFailed
	}
}
