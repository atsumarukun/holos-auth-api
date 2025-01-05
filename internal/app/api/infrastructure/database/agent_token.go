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
	ErrRequiredAgentToken = status.Error(http.StatusInternalServerError, "agent token is required")
)

type agentTokenDBRepository struct {
	db *sqlx.DB
}

func NewAgentTokenDBRepository(db *sqlx.DB) repository.AgentTokenRepository {
	return &agentTokenDBRepository{
		db: db,
	}
}

func (r *agentTokenDBRepository) Save(ctx context.Context, agentToken *entity.AgentToken) error {
	if agentToken == nil {
		return ErrRequiredAgentToken
	}

	driver := getDriver(ctx, r.db)
	agentTokenModel := transformer.ToAgentTokenModel(agentToken)

	_, err := driver.NamedExecContext(
		ctx,
		`REPLACE agent_tokens (agent_id, token) VALUES (:agent_id, :token);`,
		agentTokenModel,
	)

	return err
}

func (r *agentTokenDBRepository) Delete(ctx context.Context, agentToken *entity.AgentToken) error {
	if agentToken == nil {
		return ErrRequiredAgentToken
	}

	driver := getDriver(ctx, r.db)
	agentTokenModel := transformer.ToAgentTokenModel(agentToken)

	_, err := driver.NamedExecContext(
		ctx,
		`DELETE FROM agent_tokens WHERE agent_id = :agent_id;`,
		agentTokenModel,
	)

	return err
}

func (r *agentTokenDBRepository) FindOneByAgentIDAndUserID(ctx context.Context, agentID uuid.UUID, userID uuid.UUID) (*entity.AgentToken, error) {
	var agentToken model.AgentTokenModel
	driver := getDriver(ctx, r.db)

	if err := driver.QueryRowxContext(
		ctx,
		`SELECT
			agent_tokens.agent_id,
			agent_tokens.token
		FROM
			agent_tokens
			INNER JOIN agents ON agent_tokens.agent_id = agents.id
		WHERE
			agent_tokens.agent_id = ?
			AND agents.user_id = ?
		LIMIT 1;`,
		agentID,
		userID,
	).StructScan(&agentToken); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return transformer.ToAgentTokenEntity(&agentToken), nil
}
