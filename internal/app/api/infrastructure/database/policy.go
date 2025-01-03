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
	ErrRequiredPolicy = status.Error(http.StatusInternalServerError, "policy is required")
)

type policyDBRepository struct {
	db *sqlx.DB
}

func NewPolicyDBRepository(db *sqlx.DB) repository.PolicyRepository {
	return &policyDBRepository{
		db: db,
	}
}

func (r *policyDBRepository) Create(ctx context.Context, policy *entity.Policy) error {
	if policy == nil {
		return ErrRequiredPolicy
	}

	driver := getSqlxDriver(ctx, r.db)
	policyModel, err := transformer.ToPolicyModel(policy)
	if err != nil {
		return err
	}

	_, err = driver.NamedExecContext(
		ctx,
		`INSERT INTO policies (id, user_id, name, service, path, methods, created_at, updated_at) VALUES (:id, :user_id, :name, :service, :path, :methods, :created_at, :updated_at);`,
		policyModel,
	)

	return err
}

func (r *policyDBRepository) Update(ctx context.Context, policy *entity.Policy) error {
	if policy == nil {
		return ErrRequiredPolicy
	}

	driver := getSqlxDriver(ctx, r.db)
	policyModel, err := transformer.ToPolicyModel(policy)
	if err != nil {
		return err
	}

	if _, err := driver.NamedExecContext(
		ctx,
		`UPDATE policies SET user_id = :user_id, name = :name, service = :service, path = :path, methods = :methods, updated_at = :updated_at WHERE id = :id AND deleted_at IS NULL LIMIT 1;`,
		policyModel,
	); err != nil {
		return err
	}

	return r.updateAgents(ctx, policy.ID, policy.Agents)
}

func (r *policyDBRepository) Delete(ctx context.Context, policy *entity.Policy) error {
	if policy == nil {
		return ErrRequiredPolicy
	}

	driver := getSqlxDriver(ctx, r.db)
	policyModel, err := transformer.ToPolicyModel(policy)
	if err != nil {
		return err
	}

	_, err = driver.NamedExecContext(
		ctx,
		`UPDATE policies SET updated_at = updated_at, deleted_at = NOW(6) WHERE id = :id AND deleted_at IS NULL LIMIT 1;`,
		policyModel,
	)

	return err
}

func (r *policyDBRepository) FindOneByIDAndUserIDAndNotDeleted(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.Policy, error) {
	var policy model.PolicyModel
	driver := getSqlxDriver(ctx, r.db)

	if err := driver.QueryRowxContext(
		ctx,
		`SELECT
			policies.id,
			policies.user_id,
			policies.name,
			policies.service,
			policies.path,
			policies.methods,
			policies.created_at,
			policies.updated_at,
			GROUP_CONCAT(permissions.agent_id ORDER BY permissions.agent_id) as agents
		FROM
			policies
			LEFT JOIN permissions ON policies.id = permissions.policy_id
		WHERE
			policies.id = ?
			AND policies.user_id = ?
			AND policies.deleted_at IS NULL
		GROUP BY
			policies.id
		LIMIT 1;`,
		id,
		userID,
	).StructScan(&policy); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return transformer.ToPolicyEntity(&policy)
}

func (r *policyDBRepository) FindByUserIDAndNotDeleted(ctx context.Context, userID uuid.UUID) ([]*entity.Policy, error) {
	var policies []*model.PolicyModel
	driver := getSqlxDriver(ctx, r.db)

	rows, err := driver.QueryxContext(
		ctx,
		`SELECT id, user_id, name, service, path, methods, created_at, updated_at FROM policies WHERE user_id = ? AND deleted_at IS NULL;`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var policy model.PolicyModel
		if err := rows.StructScan(&policy); err != nil {
			return nil, err
		}
		policies = append(policies, &policy)
	}

	return transformer.ToPolicyEntities(policies)
}

func (r *policyDBRepository) FindByIDsAndUserIDAndNotDeleted(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) ([]*entity.Policy, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var policies []*model.PolicyModel
	driver := getSqlxDriver(ctx, r.db)

	query, args, err := sqlx.Named(
		`SELECT id, user_id, name, service, path, methods, created_at, updated_at FROM policies WHERE id IN (:ids) AND user_id = :user_id AND deleted_at IS NULL;`,
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
		var policy model.PolicyModel
		if err := rows.StructScan(&policy); err != nil {
			return nil, err
		}
		policies = append(policies, &policy)
	}

	return transformer.ToPolicyEntities(policies)
}

func (r *policyDBRepository) updateAgents(ctx context.Context, id uuid.UUID, agentIDs []uuid.UUID) error {
	driver := getSqlxDriver(ctx, r.db)

	if _, err := driver.NamedExecContext(
		ctx,
		`DELETE FROM permissions WHERE policy_id = :policy_id;`,
		map[string]interface{}{"policy_id": id},
	); err != nil {
		return err
	}

	if len(agentIDs) == 0 {
		return nil
	}

	args := make([]map[string]interface{}, len(agentIDs))
	for i, agentID := range agentIDs {
		args[i] = map[string]interface{}{
			"agent_id":  agentID,
			"policy_id": id,
		}
	}
	if _, err := driver.NamedExecContext(
		ctx,
		`INSERT INTO permissions (agent_id, policy_id) VALUES (:agent_id, :policy_id);`,
		args,
	); err != nil {
		return err
	}

	return nil
}
