//go:generate mockgen -source=$GOFILE -destination=../../../../test/mock/usecase/$GOFILE
package usecase

import (
	"context"
	"holos-auth-api/internal/app/api/domain"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/repository"
	"holos-auth-api/internal/app/api/pkg/apierr"
	"holos-auth-api/internal/app/api/usecase/dto"
	"net/http"

	"github.com/google/uuid"
)

var (
	ErrPolicyNotFound = apierr.NewApiError(http.StatusNotFound, "policy not found")
)

type PolicyUsecase interface {
	Create(context.Context, uuid.UUID, string, string, string, []string) (*dto.PolicyDTO, apierr.ApiError)
	Update(context.Context, uuid.UUID, uuid.UUID, string, string, string, []string) (*dto.PolicyDTO, apierr.ApiError)
	Delete(context.Context, uuid.UUID, uuid.UUID) apierr.ApiError
}

type policyUsecase struct {
	transactionObject domain.TransactionObject
	policyRepository  repository.PolicyRepository
}

func NewPolicyUsecase(transactionObject domain.TransactionObject, policyRepository repository.PolicyRepository) PolicyUsecase {
	return &policyUsecase{
		transactionObject: transactionObject,
		policyRepository:  policyRepository,
	}
}

func (u *policyUsecase) Create(ctx context.Context, userID uuid.UUID, name string, service string, path string, allowedMethods []string) (*dto.PolicyDTO, apierr.ApiError) {
	policy, err := entity.NewPolicy(userID, name, service, path, allowedMethods)
	if err != nil {
		return nil, err
	}

	if err := u.policyRepository.Create(ctx, policy); err != nil {
		return nil, err
	}

	return dto.NewPolicyDTO(policy.ID, policy.UserID, policy.Name, policy.Service, policy.Path, policy.AllowedMethods, policy.CreatedAt, policy.UpdatedAt), nil
}

func (u *policyUsecase) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, name string, service string, path string, allowedMethods []string) (*dto.PolicyDTO, apierr.ApiError) {
	var policy *entity.Policy

	if err := u.transactionObject.Transaction(ctx, func(ctx context.Context) apierr.ApiError {
		var err apierr.ApiError
		policy, err = u.policyRepository.FindOneByIDAndUserIDAndNotDeleted(ctx, id, userID)
		if err != nil {
			return err
		}
		if policy == nil {
			return ErrPolicyNotFound
		}

		if err := policy.SetName(name); err != nil {
			return err
		}
		if err := policy.SetService(service); err != nil {
			return err
		}
		if err := policy.SetPath(path); err != nil {
			return err
		}
		if err := policy.SetAllowedMethods(allowedMethods); err != nil {
			return err
		}

		return u.policyRepository.Update(ctx, policy)
	}); err != nil {
		return nil, err
	}

	return dto.NewPolicyDTO(policy.ID, policy.UserID, policy.Name, policy.Service, policy.Path, policy.AllowedMethods, policy.CreatedAt, policy.UpdatedAt), nil
}

func (u *policyUsecase) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) apierr.ApiError {
	return u.transactionObject.Transaction(ctx, func(ctx context.Context) apierr.ApiError {
		policy, err := u.policyRepository.FindOneByIDAndUserIDAndNotDeleted(ctx, id, userID)
		if err != nil {
			return err
		}
		if policy == nil {
			return ErrPolicyNotFound
		}

		return u.policyRepository.Delete(ctx, policy)
	})
}
