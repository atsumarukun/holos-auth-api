package database_test

import (
	"context"
	"database/sql"
	"errors"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/infrastructure/database"
	"holos-auth-api/test"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestAgent_Create(t *testing.T) {
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name        string
		inputAgent  *entity.Agent
		expectError error
		setMockDB   func(sqlmock.Sqlmock)
	}{
		{
			name:        "success",
			inputAgent:  agent,
			expectError: nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO agents (id, user_id, name, created_at, updated_at) VALUES (?, ?, ?, ?, ?);")).
					WithArgs(agent.ID, agent.UserID, agent.Name, agent.CreatedAt, agent.UpdatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
		},
		{
			name:        "create error",
			inputAgent:  agent,
			expectError: sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO agents (id, user_id, name, created_at, updated_at) VALUES (?, ?, ?, ?, ?);")).
					WithArgs(agent.ID, agent.UserID, agent.Name, agent.CreatedAt, agent.UpdatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:        "no agent",
			inputAgent:  nil,
			expectError: database.ErrRequiredAgent,
			setMockDB:   func(mock sqlmock.Sqlmock) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := test.NewMockDB(t)
			defer db.Close()

			ctx := context.Background()

			tt.setMockDB(mock)

			r := database.NewAgentDBRepository(db)
			if err := r.Create(ctx, tt.inputAgent); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func TestAgent_Update(t *testing.T) {
	agentWithoutPolicies, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}

	policy, err := entity.NewPolicy(agentWithoutPolicies.UserID, "name", "ALLOW", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}
	agentWithPolicies, err := entity.NewAgent(agentWithoutPolicies.UserID, "name")
	if err != nil {
		t.Error(err.Error())
	}
	agentWithPolicies.SetPolicies([]*entity.Policy{policy})

	tests := []struct {
		name        string
		inputAgent  *entity.Agent
		expectError error
		setMockDB   func(sqlmock.Sqlmock)
	}{
		{
			name:        "without permissions",
			inputAgent:  agentWithoutPolicies,
			expectError: nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE agents SET user_id = ?, name = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
					WithArgs(agentWithoutPolicies.UserID, agentWithoutPolicies.Name, agentWithoutPolicies.UpdatedAt, agentWithoutPolicies.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM permissions WHERE agent_id = ?;")).
					WithArgs(agentWithoutPolicies.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
		},
		{
			name:        "with permissions",
			inputAgent:  agentWithPolicies,
			expectError: nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE agents SET user_id = ?, name = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
					WithArgs(agentWithPolicies.UserID, agentWithPolicies.Name, agentWithPolicies.UpdatedAt, agentWithPolicies.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM permissions WHERE agent_id = ?;")).
					WithArgs(agentWithPolicies.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO permissions (agent_id, policy_id) VALUES (?, ?);")).
					WithArgs(agentWithPolicies.ID, agentWithPolicies.Policies[0]).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
		},
		{
			name:        "update agent error",
			inputAgent:  agentWithoutPolicies,
			expectError: sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE agents SET user_id = ?, name = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
					WithArgs(agentWithoutPolicies.UserID, agentWithoutPolicies.Name, agentWithoutPolicies.UpdatedAt, agentWithoutPolicies.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:        "delete permissions error",
			inputAgent:  agentWithoutPolicies,
			expectError: sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE agents SET user_id = ?, name = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
					WithArgs(agentWithoutPolicies.UserID, agentWithoutPolicies.Name, agentWithoutPolicies.UpdatedAt, agentWithoutPolicies.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM permissions WHERE agent_id = ?;")).
					WithArgs(agentWithoutPolicies.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:        "create permissions error",
			inputAgent:  agentWithPolicies,
			expectError: sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE agents SET user_id = ?, name = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
					WithArgs(agentWithPolicies.UserID, agentWithPolicies.Name, agentWithPolicies.UpdatedAt, agentWithPolicies.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM permissions WHERE agent_id = ?;")).
					WithArgs(agentWithPolicies.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO permissions (agent_id, policy_id) VALUES (?, ?);")).
					WithArgs(agentWithPolicies.ID, agentWithPolicies.Policies[0]).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:        "no agent",
			inputAgent:  nil,
			expectError: database.ErrRequiredAgent,
			setMockDB:   func(mock sqlmock.Sqlmock) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := test.NewMockDB(t)
			defer db.Close()

			ctx := context.Background()

			tt.setMockDB(mock)

			r := database.NewAgentDBRepository(db)
			if err := r.Update(ctx, tt.inputAgent); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func TestAgent_Delete(t *testing.T) {
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name        string
		inputAgent  *entity.Agent
		expectError error
		setMockDB   func(sqlmock.Sqlmock)
	}{
		{
			name:        "success",
			inputAgent:  agent,
			expectError: nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE agents SET updated_at = updated_at, deleted_at = NOW(6) WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
					WithArgs(agent.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
		},
		{
			name:        "delete error",
			inputAgent:  agent,
			expectError: sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE agents SET updated_at = updated_at, deleted_at = NOW(6) WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
					WithArgs(agent.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:        "no agent",
			inputAgent:  nil,
			expectError: database.ErrRequiredAgent,
			setMockDB:   func(mock sqlmock.Sqlmock) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := test.NewMockDB(t)
			defer db.Close()

			ctx := context.Background()

			tt.setMockDB(mock)

			r := database.NewAgentDBRepository(db)
			if err := r.Delete(ctx, tt.inputAgent); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func TestAgent_FindOneByIDAndUserIDAndNotDeleted(t *testing.T) {
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name         string
		inputID      uuid.UUID
		inputUserID  uuid.UUID
		expectResult *entity.Agent
		expectError  error
		setMockDB    func(sqlmock.Sqlmock)
	}{
		{
			name:         "found",
			inputID:      agent.ID,
			inputUserID:  agent.UserID,
			expectResult: agent,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
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
				)).
					WithArgs(agent.ID, agent.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}).
							AddRow(agent.ID, agent.UserID, agent.Name, agent.CreatedAt, agent.UpdatedAt),
					).
					WillReturnError(nil)
			},
		},
		{
			name:         "not found",
			inputID:      agent.ID,
			inputUserID:  agent.UserID,
			expectResult: nil,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
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
				)).
					WithArgs(agent.ID, agent.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}),
					).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:         "find error",
			inputID:      agent.ID,
			inputUserID:  agent.UserID,
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
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
				)).
					WithArgs(agent.ID, agent.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}),
					).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := test.NewMockDB(t)
			defer db.Close()

			ctx := context.Background()

			tt.setMockDB(mock)

			r := database.NewAgentDBRepository(db)
			result, err := r.FindOneByIDAndUserIDAndNotDeleted(ctx, tt.inputID, tt.inputUserID)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func TestAgent_FindOneByTokenAndNotDeleted(t *testing.T) {
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}
	agentToken, err := entity.NewAgentToken(agent.ID)
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name         string
		inputToken   string
		expectResult *entity.Agent
		expectError  error
		setMockDB    func(sqlmock.Sqlmock)
	}{
		{
			name:         "found",
			inputToken:   agentToken.Token,
			expectResult: agent,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
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
				)).
					WithArgs(agentToken.Token).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}).
							AddRow(agent.ID, agent.UserID, agent.Name, agent.CreatedAt, agent.UpdatedAt),
					).
					WillReturnError(nil)
			},
		},
		{
			name:         "not found",
			inputToken:   agentToken.Token,
			expectResult: nil,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
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
				)).
					WithArgs(agentToken.Token).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}),
					).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:         "find error",
			inputToken:   agentToken.Token,
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
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
				)).
					WithArgs(agentToken.Token).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}),
					).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := test.NewMockDB(t)
			defer db.Close()

			ctx := context.Background()

			tt.setMockDB(mock)

			r := database.NewAgentDBRepository(db)
			result, err := r.FindOneByTokenAndNotDeleted(ctx, tt.inputToken)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func TestAgent_FindByNamePrefixAndUserIDAndNotDeleted(t *testing.T) {
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name         string
		inputKeyword string
		inputUserID  uuid.UUID
		expectResult []*entity.Agent
		expectError  error
		setMockDB    func(sqlmock.Sqlmock)
	}{
		{
			name:         "found",
			inputKeyword: "name",
			inputUserID:  agent.UserID,
			expectResult: []*entity.Agent{agent},
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, created_at, updated_at FROM agents WHERE name LIKE ? AND user_id = ? AND deleted_at IS NULL;")).
					WithArgs("name%", agent.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}).
							AddRow(agent.ID, agent.UserID, agent.Name, agent.CreatedAt, agent.UpdatedAt),
					).
					WillReturnError(nil)
			},
		},
		{
			name:         "not found",
			inputKeyword: "keyword",
			inputUserID:  agent.UserID,
			expectResult: []*entity.Agent{},
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, created_at, updated_at FROM agents WHERE name LIKE ? AND user_id = ? AND deleted_at IS NULL;")).
					WithArgs("keyword%", agent.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}),
					).
					WillReturnError(nil)
			},
		},
		{
			name:         "find error",
			inputKeyword: "name",
			inputUserID:  agent.UserID,
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, created_at, updated_at FROM agents WHERE name LIKE ? AND user_id = ? AND deleted_at IS NULL;")).
					WithArgs("name%", agent.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}),
					).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := test.NewMockDB(t)
			defer db.Close()

			ctx := context.Background()

			tt.setMockDB(mock)

			r := database.NewAgentDBRepository(db)
			result, err := r.FindByNamePrefixAndUserIDAndNotDeleted(ctx, tt.inputKeyword, tt.inputUserID)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func TestAgent_FindByIDsAndUserIDAndNotDeleted(t *testing.T) {
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name         string
		inputIDs     []uuid.UUID
		inputUserID  uuid.UUID
		expectResult []*entity.Agent
		expectError  error
		setMockDB    func(sqlmock.Sqlmock)
	}{
		{
			name:         "found",
			inputIDs:     []uuid.UUID{agent.ID},
			inputUserID:  agent.UserID,
			expectResult: []*entity.Agent{agent},
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, created_at, updated_at FROM agents WHERE id IN (?) AND user_id = ? AND deleted_at IS NULL;")).
					WithArgs(agent.ID, agent.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}).
							AddRow(agent.ID, agent.UserID, agent.Name, agent.CreatedAt, agent.UpdatedAt),
					).
					WillReturnError(nil)
			},
		},
		{
			name:         "not found",
			inputIDs:     []uuid.UUID{agent.ID},
			inputUserID:  agent.UserID,
			expectResult: []*entity.Agent{},
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, created_at, updated_at FROM agents WHERE id IN (?) AND user_id = ? AND deleted_at IS NULL;")).
					WithArgs(agent.ID, agent.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}),
					).
					WillReturnError(nil)
			},
		},
		{
			name:         "find error",
			inputIDs:     []uuid.UUID{agent.ID},
			inputUserID:  agent.UserID,
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, created_at, updated_at FROM agents WHERE id IN (?) AND user_id = ? AND deleted_at IS NULL;")).
					WithArgs(agent.ID, agent.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}),
					).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := test.NewMockDB(t)
			defer db.Close()

			ctx := context.Background()

			tt.setMockDB(mock)

			r := database.NewAgentDBRepository(db)
			result, err := r.FindByIDsAndUserIDAndNotDeleted(ctx, tt.inputIDs, tt.inputUserID)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func TestAgent_FindByIDsAndNamePrefixAndUserIDAndNotDeleted(t *testing.T) {
	agent, err := entity.NewAgent(uuid.New(), "name")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name         string
		inputIDs     []uuid.UUID
		inputKeyword string
		inputUserID  uuid.UUID
		expectResult []*entity.Agent
		expectError  error
		setMockDB    func(sqlmock.Sqlmock)
	}{
		{
			name:         "found",
			inputIDs:     []uuid.UUID{agent.ID},
			inputKeyword: "name",
			inputUserID:  agent.UserID,
			expectResult: []*entity.Agent{agent},
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, created_at, updated_at FROM agents WHERE id IN (?) AND name LIKE ? AND user_id = ? AND deleted_at IS NULL;")).
					WithArgs(agent.ID, "name%", agent.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}).
							AddRow(agent.ID, agent.UserID, agent.Name, agent.CreatedAt, agent.UpdatedAt),
					).
					WillReturnError(nil)
			},
		},
		{
			name:         "not found",
			inputIDs:     []uuid.UUID{agent.ID},
			inputKeyword: "keyword",
			inputUserID:  agent.UserID,
			expectResult: []*entity.Agent{},
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, created_at, updated_at FROM agents WHERE id IN (?) AND name LIKE ? AND user_id = ? AND deleted_at IS NULL;")).
					WithArgs(agent.ID, "keyword%", agent.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}),
					).
					WillReturnError(nil)
			},
		},
		{
			name:         "find error",
			inputIDs:     []uuid.UUID{agent.ID},
			inputKeyword: "name",
			inputUserID:  agent.UserID,
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, created_at, updated_at FROM agents WHERE id IN (?) AND name LIKE ? AND user_id = ? AND deleted_at IS NULL;")).
					WithArgs(agent.ID, "name%", agent.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}),
					).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := test.NewMockDB(t)
			defer db.Close()

			ctx := context.Background()

			tt.setMockDB(mock)

			r := database.NewAgentDBRepository(db)
			result, err := r.FindByIDsAndNamePrefixAndUserIDAndNotDeleted(ctx, tt.inputIDs, tt.inputKeyword, tt.inputUserID)
			if !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}
			if diff := cmp.Diff(result, tt.expectResult); diff != "" {
				t.Error(diff)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err.Error())
			}
		})
	}
}
