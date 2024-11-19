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

type agentInfrastructure struct {
	db *sqlx.DB
}

func NewAgentInfrastructure(db *sqlx.DB) repository.AgentRepository {
	return &agentInfrastructure{
		db: db,
	}
}

func (ai *agentInfrastructure) Create(ctx context.Context, agent *entity.Agent) apierr.ApiError {
	driver := getSqlxDriver(ctx, ai.db)
	if _, err := driver.NamedExecContext(
		ctx,
		`INSERT INTO agents (id, user_id, name, created_at, updated_at) VALUES (:id, :user_id, :name, :created_at, :updated_at);`,
		agent,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (ai *agentInfrastructure) Update(ctx context.Context, agent *entity.Agent) apierr.ApiError {
	driver := getSqlxDriver(ctx, ai.db)
	if _, err := driver.NamedExecContext(
		ctx,
		`UPDATE agents SET user_id = :user_id, name = :name, updated_at = :updated_at WHERE id = :id AND deleted_at IS NULL LIMIT 1;`,
		agent,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (ai *agentInfrastructure) Delete(ctx context.Context, agent *entity.Agent) apierr.ApiError {
	driver := getSqlxDriver(ctx, ai.db)
	if _, err := driver.NamedExecContext(
		ctx,
		`UPDATE agents SET updated_at = updated_at, deleted_at = NOW(6) WHERE id = :id AND deleted_at IS NULL LIMIT 1;`,
		agent,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (ai *agentInfrastructure) FindOneByIDAndUserIDAndNotDeleted(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Agent, apierr.ApiError) {
	var agent entity.Agent
	driver := getSqlxDriver(ctx, ai.db)
	if err := driver.QueryRowxContext(
		ctx,
		`SELECT id, user_id, name, created_at, updated_at FROM agents WHERE id = ? AND user_id = ? AND deleted_at IS NULL LIMIT 1;`,
		id,
		userID,
	).StructScan(&agent); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, apierr.NewApiError(http.StatusInternalServerError, err.Error())
		}
	}
	return &agent, nil
}

func (ai *agentInfrastructure) FindOneByUserIDAndName(ctx context.Context, userID uuid.UUID, name string) (*entity.Agent, apierr.ApiError) {
	var agent entity.Agent
	driver := getSqlxDriver(ctx, ai.db)
	if err := driver.QueryRowxContext(
		ctx,
		`SELECT id, user_id, name, created_at, updated_at FROM agents WHERE user_id = ? AND name = ? LIMIT 1;`,
		userID,
		name,
	).StructScan(&agent); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, apierr.NewApiError(http.StatusInternalServerError, err.Error())
		}
	}
	return &agent, nil
}
