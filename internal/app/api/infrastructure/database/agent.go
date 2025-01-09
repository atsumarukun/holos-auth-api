package database

import (
	"context"
	"database/sql"
	"errors"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/repository"
	"holos-auth-api/internal/app/api/infrastructure/model"
	"holos-auth-api/internal/app/api/infrastructure/transformer"
	"holos-auth-api/internal/app/api/pkg/status"
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrRequiredAgent = status.Error(http.StatusInternalServerError, "agent is required")
)

type agentDBRepository struct {
	db *sqlx.DB
}

func NewAgentDBRepository(db *sqlx.DB) repository.AgentRepository {
	return &agentDBRepository{
		db: db,
	}
}

func (r *agentDBRepository) Create(ctx context.Context, agent *entity.Agent) error {
	if agent == nil {
		return ErrRequiredAgent
	}

	driver := getDriver(ctx, r.db)
	agentModel := transformer.ToAgentModel(agent)

	_, err := driver.NamedExecContext(
		ctx,
		`INSERT INTO agents (id, user_id, name, created_at, updated_at) VALUES (:id, :user_id, :name, :created_at, :updated_at);`,
		agentModel,
	)

	return err
}

func (r *agentDBRepository) Update(ctx context.Context, agent *entity.Agent) error {
	if agent == nil {
		return ErrRequiredAgent
	}

	driver := getDriver(ctx, r.db)
	agentModel := transformer.ToAgentModel(agent)

	if _, err := driver.NamedExecContext(
		ctx,
		`UPDATE agents SET user_id = :user_id, name = :name, updated_at = :updated_at WHERE id = :id AND deleted_at IS NULL LIMIT 1;`,
		agentModel,
	); err != nil {
		return err
	}

	return r.updatePolicies(ctx, agent.ID, agent.Policies)
}

func (r *agentDBRepository) Delete(ctx context.Context, agent *entity.Agent) error {
	if agent == nil {
		return ErrRequiredAgent
	}

	driver := getDriver(ctx, r.db)
	agentModel := transformer.ToAgentModel(agent)

	_, err := driver.NamedExecContext(
		ctx,
		`UPDATE agents SET updated_at = updated_at, deleted_at = NOW(6) WHERE id = :id AND deleted_at IS NULL LIMIT 1;`,
		agentModel,
	)

	return err
}

func (r *agentDBRepository) FindOneByIDAndUserIDAndNotDeleted(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Agent, error) {
	var agent model.AgentModel
	driver := getDriver(ctx, r.db)

	if err := driver.QueryRowxContext(
		ctx,
		`SELECT
			agents.id,
			agents.user_id,
			agents.name,
			agents.created_at,
			agents.updated_at,
			GROUP_CONCAT(permissions.policy_id ORDER BY permissions.policy_id) as policies
		FROM
			agents
			LEFT JOIN permissions ON agents.id = permissions.agent_id
		WHERE
			agents.id = ?
			AND agents.user_id = ?
			AND agents.deleted_at IS NULL
		GROUP BY
			agents.id
		LIMIT 1;`,
		id,
		userID,
	).StructScan(&agent); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return transformer.ToAgentEntity(&agent)
}

func (r *agentDBRepository) FindOneByTokenAndNotDeleted(ctx context.Context, token string) (*entity.Agent, error) {
	var agent model.AgentModel
	driver := getDriver(ctx, r.db)

	if err := driver.QueryRowxContext(
		ctx,
		`SELECT
			agents.id,
			agents.user_id,
			agents.name,
			agents.created_at,
			agents.updated_at,
			GROUP_CONCAT(permissions.policy_id ORDER BY permissions.policy_id) as policies
		FROM
			agents
			INNER JOIN agent_tokens ON agents.id = agent_tokens.agent_id
			LEFT JOIN permissions ON agents.id = permissions.agent_id
		WHERE
			agent_tokens.token = ?
			AND agents.deleted_at IS NULL
		GROUP BY
			agents.id
		LIMIT 1;`,
		token,
	).StructScan(&agent); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return transformer.ToAgentEntity(&agent)
}

func (r *agentDBRepository) FindByUserIDAndNotDeleted(ctx context.Context, userID uuid.UUID) ([]*entity.Agent, error) {
	agents := []*model.AgentModel{}
	driver := getDriver(ctx, r.db)

	rows, err := driver.QueryxContext(
		ctx,
		`SELECT id, user_id, name, created_at, updated_at FROM agents WHERE user_id = ? AND deleted_at IS NULL;`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var agent model.AgentModel
		if err := rows.StructScan(&agent); err != nil {
			return nil, err
		}
		agents = append(agents, &agent)
	}

	return transformer.ToAgentEntities(agents)
}

func (r *agentDBRepository) FindByIDsAndUserIDAndNotDeleted(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) ([]*entity.Agent, error) {
	if len(ids) == 0 {
		return []*entity.Agent{}, nil
	}

	agents := []*model.AgentModel{}
	driver := getDriver(ctx, r.db)

	query, args, err := sqlx.Named(
		`SELECT id, user_id, name, created_at, updated_at FROM agents WHERE id IN (:ids) AND user_id = :user_id AND deleted_at IS NULL;`,
		map[string]interface{}{
			"ids":     ids,
			"user_id": userID,
		},
	)
	if err != nil {
		return nil, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}
	query = driver.Rebind(query)

	rows, err := driver.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var agent model.AgentModel
		if err := rows.StructScan(&agent); err != nil {
			return nil, err
		}
		agents = append(agents, &agent)
	}

	return transformer.ToAgentEntities(agents)
}

func (r *agentDBRepository) updatePolicies(ctx context.Context, id uuid.UUID, policieIDs []uuid.UUID) error {
	driver := getDriver(ctx, r.db)

	if _, err := driver.NamedExecContext(
		ctx,
		`DELETE FROM permissions WHERE agent_id = :agent_id;`,
		map[string]interface{}{"agent_id": id},
	); err != nil {
		return err
	}

	if len(policieIDs) == 0 {
		return nil
	}

	args := make([]map[string]interface{}, len(policieIDs))
	for i, policyID := range policieIDs {
		args[i] = map[string]interface{}{
			"agent_id":  id,
			"policy_id": policyID,
		}
	}
	_, err := driver.NamedExecContext(
		ctx,
		`INSERT INTO permissions (agent_id, policy_id) VALUES (:agent_id, :policy_id);`,
		args,
	)

	return err
}
