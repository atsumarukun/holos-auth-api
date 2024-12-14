package database

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

type agentDBRepository struct {
	db *sqlx.DB
}

func NewAgentDBRepository(db *sqlx.DB) repository.AgentRepository {
	return &agentDBRepository{
		db: db,
	}
}

func (r *agentDBRepository) Create(ctx context.Context, agent *entity.Agent) apierr.ApiError {
	driver := getSqlxDriver(ctx, r.db)
	agentModel := r.convertToModel(agent)
	if _, err := driver.NamedExecContext(
		ctx,
		`INSERT INTO agents (id, user_id, name, created_at, updated_at) VALUES (:id, :user_id, :name, :created_at, :updated_at);`,
		agentModel,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *agentDBRepository) Update(ctx context.Context, agent *entity.Agent) apierr.ApiError {
	driver := getSqlxDriver(ctx, r.db)
	agentModel := r.convertToModel(agent)
	if _, err := driver.NamedExecContext(
		ctx,
		`UPDATE agents SET user_id = :user_id, name = :name, updated_at = :updated_at WHERE id = :id AND deleted_at IS NULL LIMIT 1;`,
		agentModel,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *agentDBRepository) Delete(ctx context.Context, agent *entity.Agent) apierr.ApiError {
	driver := getSqlxDriver(ctx, r.db)
	agentModel := r.convertToModel(agent)
	if _, err := driver.NamedExecContext(
		ctx,
		`UPDATE agents SET updated_at = updated_at, deleted_at = NOW(6) WHERE id = :id AND deleted_at IS NULL LIMIT 1;`,
		agentModel,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *agentDBRepository) FindOneByIDAndUserIDAndNotDeleted(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Agent, apierr.ApiError) {
	var agent model.AgentModel
	driver := getSqlxDriver(ctx, r.db)
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
	return r.convertToEntity(&agent), nil
}

func (r *agentDBRepository) FindByUserIDAndNotDeleted(ctx context.Context, userID uuid.UUID) ([]*entity.Agent, apierr.ApiError) {
	var agents []*model.AgentModel
	driver := getSqlxDriver(ctx, r.db)
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
		var agent model.AgentModel
		if err := rows.StructScan(&agent); err != nil {
			return nil, apierr.NewApiError(http.StatusInternalServerError, err.Error())
		}
		agents = append(agents, &agent)
	}
	return r.convertToEntities(agents), nil
}

func (r *agentDBRepository) FindByIDsAndUserIDAndNotDeleted(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) ([]*entity.Agent, apierr.ApiError) {
	var agents []*model.AgentModel
	driver := getSqlxDriver(ctx, r.db)

	query, args, err := sqlx.Named(
		`SELECT id, user_id, name, created_at, updated_at FROM agents WHERE id IN (:ids) AND user_id = :user_id AND deleted_at IS NULL;`,
		map[string]interface{}{
			"ids":     ids,
			"user_id": userID,
		},
	)
	if err != nil {
		return nil, apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	query = driver.Rebind(query)

	rows, err := driver.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var agent model.AgentModel
		if err := rows.StructScan(&agent); err != nil {
			return nil, apierr.NewApiError(http.StatusInternalServerError, err.Error())
		}
		agents = append(agents, &agent)
	}
	return r.convertToEntities(agents), nil
}

func (r *agentDBRepository) UpdatePolicies(ctx context.Context, id uuid.UUID, policies []*entity.Policy) apierr.ApiError {
	driver := getSqlxDriver(ctx, r.db)

	if _, err := driver.NamedExecContext(
		ctx,
		`DELETE FROM permissions WHERE agent_id = :agent_id;`,
		map[string]interface{}{"agent_id": id},
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}

	if len(policies) == 0 {
		return nil
	}

	args := make([]map[string]interface{}, len(policies))
	for i, policy := range policies {
		args[i] = map[string]interface{}{
			"agent_id":  id,
			"policy_id": policy.ID,
		}
	}
	if _, err := driver.NamedExecContext(
		ctx,
		`INSERT INTO permissions (agent_id, policy_id) VALUES (:agent_id, :policy_id);`,
		args,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (r *agentDBRepository) GetPolicies(ctx context.Context, id uuid.UUID, userID uuid.UUID) ([]*entity.Policy, apierr.ApiError) {
	var policies []*model.PolicyModel
	driver := getSqlxDriver(ctx, r.db)
	rows, err := driver.QueryxContext(
		ctx,
		`SELECT
			policies.id,
			policies.user_id,
			policies.name,
			policies.service,
			policies.path,
			policies.methods,
			policies.created_at,
			policies.updated_at
		FROM
			policies
			LEFT JOIN permissions ON policies.id = permissions.policy_id
		WHERE
			policies.user_id = ?
			AND permissions.agent_id = ?
			AND policies.deleted_at IS NULL;`,
		userID,
		id,
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
	entities := make([]*entity.Policy, len(policies))
	for i, policy := range policies {
		var methods []string
		if err := json.Unmarshal([]byte(policy.Methods), &methods); err != nil {
			return nil, apierr.NewApiError(http.StatusInternalServerError, err.Error())
		}
		entities[i] = entity.RestorePolicy(policy.ID, policy.UserID, policy.Name, policy.Service, policy.Path, methods, policy.CreatedAt, policy.UpdatedAt)
	}
	return entities, nil
}

func (r *agentDBRepository) convertToModel(agent *entity.Agent) *model.AgentModel {
	return model.NewAgentModel(agent.ID, agent.UserID, agent.Name, agent.CreatedAt, agent.UpdatedAt)
}

func (r *agentDBRepository) convertToEntity(agent *model.AgentModel) *entity.Agent {
	return entity.RestoreAgent(agent.ID, agent.UserID, agent.Name, agent.CreatedAt, agent.UpdatedAt)
}

func (r *agentDBRepository) convertToEntities(agents []*model.AgentModel) []*entity.Agent {
	entities := make([]*entity.Agent, len(agents))
	for i, agent := range agents {
		entities[i] = r.convertToEntity(agent)
	}
	return entities
}
