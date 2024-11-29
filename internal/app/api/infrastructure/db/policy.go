package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/repository"
	"holos-auth-api/internal/app/api/infrastructure/model"
	"holos-auth-api/internal/app/api/pkg/apierr"
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type policyInfrastructure struct {
	db *sqlx.DB
}

func NewPolicyInfrastructure(db *sqlx.DB) repository.PolicyRepository {
	return &policyInfrastructure{
		db: db,
	}
}

func (i *policyInfrastructure) Create(ctx context.Context, policy *entity.Policy) apierr.ApiError {
	driver := getSqlxDriver(ctx, i.db)
	policyModel, err := i.convertToModel(policy)
	if err != nil {
		return err
	}
	if _, err := driver.NamedExecContext(
		ctx,
		`INSERT INTO policies (id, user_id, name, service, path, allowed_methods, created_at, updated_at) VALUES (:id, :user_id, :name, :service, :path, :allowed_methods, :created_at, :updated_at);`,
		policyModel,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (i *policyInfrastructure) Update(ctx context.Context, policy *entity.Policy) apierr.ApiError {
	driver := getSqlxDriver(ctx, i.db)
	policyModel, err := i.convertToModel(policy)
	if err != nil {
		return err
	}
	if _, err := driver.NamedExecContext(
		ctx,
		`UPDATE policies SET user_id = :user_id, name = :name, service = :service, path = :path, allowed_methods = :allowed_methods, updated_at = :updated_at WHERE id = :id AND deleted_at IS NULL LIMIT 1;`,
		policyModel,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (i *policyInfrastructure) Delete(ctx context.Context, policy *entity.Policy) apierr.ApiError {
	driver := getSqlxDriver(ctx, i.db)
	policyModel, err := i.convertToModel(policy)
	if err != nil {
		return err
	}
	if _, err := driver.NamedExecContext(
		ctx,
		`UPDATE policies SET updated_at = updated_at, deleted_at = NOW(6) WHERE id = :id AND deleted_at IS NULL LIMIT 1;`,
		policyModel,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (i *policyInfrastructure) FindOneByIDAndUserIDAndNotDeleted(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Policy, apierr.ApiError) {
	var policy model.PolicyModel
	driver := getSqlxDriver(ctx, i.db)
	if err := driver.QueryRowxContext(
		ctx,
		`SELECT id, user_id, name, service, path, allowed_methods, created_at, updated_at FROM policies WHERE id = ? AND user_id = ? AND deleted_at IS NULL LIMIT 1;`,
		id,
		userID,
	).StructScan(&policy); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, apierr.NewApiError(http.StatusInternalServerError, err.Error())
		}
	}
	return i.convertToEntity(&policy)
}

func (i *policyInfrastructure) FindByUserIDAndNotDeleted(ctx context.Context, userID uuid.UUID) ([]*entity.Policy, apierr.ApiError) {
	var policies []*model.PolicyModel
	driver := getSqlxDriver(ctx, i.db)
	rows, err := driver.QueryxContext(
		ctx,
		`SELECT id, user_id, name, service, path, allowed_methods, created_at, updated_at FROM policies WHERE user_id = ? AND deleted_at IS NULL;`,
		userID,
	)
	if err != nil {
		return nil, apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var policy model.PolicyModel
		if err := rows.StructScan(&policy); err != nil {
			return nil, apierr.NewApiError(http.StatusInternalServerError, err.Error())
		}
		policies = append(policies, &policy)
	}
	return i.convertToEntities(policies)
}

func (i *policyInfrastructure) convertToModel(policy *entity.Policy) (*model.PolicyModel, apierr.ApiError) {
	allowedMethods, err := json.Marshal(policy.AllowedMethods)
	if err != nil {
		return nil, apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return model.NewPolicyModel(policy.ID, policy.UserID, policy.Name, policy.Service, policy.Path, string(allowedMethods), policy.CreatedAt, policy.UpdatedAt), nil
}

func (i *policyInfrastructure) convertToEntity(policy *model.PolicyModel) (*entity.Policy, apierr.ApiError) {
	var allowedMethods []string
	if err := json.Unmarshal([]byte(policy.AllowedMethods), &allowedMethods); err != nil {
		return nil, apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return entity.RestorePolicy(policy.ID, policy.UserID, policy.Name, policy.Service, policy.Path, allowedMethods, policy.CreatedAt, policy.UpdatedAt), nil
}

func (i *policyInfrastructure) convertToEntities(policies []*model.PolicyModel) ([]*entity.Policy, apierr.ApiError) {
	entities := make([]*entity.Policy, len(policies))
	var err apierr.ApiError
	for idx, policy := range policies {
		entities[idx], err = i.convertToEntity(policy)
		if err != nil {
			return nil, err
		}
	}
	return entities, nil
}
