package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/repository"
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
	if _, err := driver.NamedExecContext(
		ctx,
		`INSERT INTO policies (id, user_id, name, service, path, allowed_methods, created_at, updated_at) VALUES (:id, :user_id, :name, :service, :path, :allowed_methods, :created_at, :updated_at);`,
		policy,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (i *policyInfrastructure) Update(ctx context.Context, policy *entity.Policy) apierr.ApiError {
	driver := getSqlxDriver(ctx, i.db)
	if _, err := driver.NamedExecContext(
		ctx,
		`UPDATE policies SET user_id = :user_id, name = :name, service = :service, path = :path, allowed_methods = :allowed_methods, updated_at = :updated_at WHERE id = :id AND deleted_at IS NULL LIMIT 1;`,
		policy,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (i *policyInfrastructure) Delete(ctx context.Context, policy *entity.Policy) apierr.ApiError {
	driver := getSqlxDriver(ctx, i.db)
	if _, err := driver.NamedExecContext(
		ctx,
		`UPDATE policies SET updated_at = updated_at, deleted_at = NOW(6) WHERE id = :id AND deleted_at IS NULL LIMIT 1;`,
		policy,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (i *policyInfrastructure) FindOneByIDAndUserIDAndNotDeleted(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Policy, apierr.ApiError) {
	var policy entity.Policy
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
	return &policy, nil
}
