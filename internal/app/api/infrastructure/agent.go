package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/repository"
	"holos-auth-api/internal/app/api/infrastructure/model"
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

func (i *agentInfrastructure) Create(ctx context.Context, agent *entity.Agent) apierr.ApiError {
	driver := getSqlxDriver(ctx, i.db)
	agentModel := i.convertToModel(agent)
	if _, err := driver.NamedExecContext(
		ctx,
		`INSERT INTO agents (id, user_id, name, created_at, updated_at) VALUES (:id, :user_id, :name, :created_at, :updated_at);`,
		agentModel,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (i *agentInfrastructure) Update(ctx context.Context, agent *entity.Agent) apierr.ApiError {
	driver := getSqlxDriver(ctx, i.db)
	agentModel := i.convertToModel(agent)
	if _, err := driver.NamedExecContext(
		ctx,
		`UPDATE agents SET user_id = :user_id, name = :name, updated_at = :updated_at WHERE id = :id AND deleted_at IS NULL LIMIT 1;`,
		agentModel,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (i *agentInfrastructure) Delete(ctx context.Context, agent *entity.Agent) apierr.ApiError {
	driver := getSqlxDriver(ctx, i.db)
	agentModel := i.convertToModel(agent)
	if _, err := driver.NamedExecContext(
		ctx,
		`UPDATE agents SET updated_at = updated_at, deleted_at = NOW(6) WHERE id = :id AND deleted_at IS NULL LIMIT 1;`,
		agentModel,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (i *agentInfrastructure) FindOneByIDAndUserIDAndNotDeleted(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Agent, apierr.ApiError) {
	var agent model.AgentModel
	driver := getSqlxDriver(ctx, i.db)
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
	return i.convertToEntity(&agent), nil
}

func (i *agentInfrastructure) FindByUserIDAndNotDeleted(ctx context.Context, userID uuid.UUID) ([]*entity.Agent, apierr.ApiError) {
	var agents []*model.AgentModel
	driver := getSqlxDriver(ctx, i.db)
	rows, err := driver.QueryxContext(
		ctx,
		`SELECT id, user_id, name, created_at, updated_at FROM agents WHERE user_id = ? AND deleted_at IS NULL;`,
		userID,
	)
	if err != nil {
		return nil, apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var agent *model.AgentModel
		if err := rows.StructScan(&agent); err != nil {
			return nil, apierr.NewApiError(http.StatusInternalServerError, err.Error())
		}
		agents = append(agents, agent)
	}
	return i.convertToEntities(agents), nil
}

func (i *agentInfrastructure) convertToModel(agent *entity.Agent) *model.AgentModel {
	return model.NewAgentModel(agent.ID, agent.UserID, agent.Name, agent.CreatedAt, agent.UpdatedAt)
}

func (i *agentInfrastructure) convertToEntity(agent *model.AgentModel) *entity.Agent {
	return entity.RestoreAgent(agent.ID, agent.UserID, agent.Name, agent.CreatedAt, agent.UpdatedAt)
}

func (i *agentInfrastructure) convertToEntities(agents []*model.AgentModel) []*entity.Agent {
	entities := make([]*entity.Agent, len(agents))
	for idx, agent := range agents {
		entities[idx] = i.convertToEntity(agent)
	}
	return entities
}
