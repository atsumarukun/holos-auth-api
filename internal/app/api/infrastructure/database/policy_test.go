package database_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/infrastructure/database"
	"holos-auth-api/test"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestPolicy_Create(t *testing.T) {
	policy, err := entity.NewPolicy(uuid.New(), "name", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name        string
		inputPolicy *entity.Policy
		expectError error
		setMockDB   func(sqlmock.Sqlmock)
	}{
		{
			name:        "success",
			inputPolicy: policy,
			expectError: nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO policies (id, user_id, name, service, path, methods, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?);")).
					WithArgs(policy.ID, policy.UserID, policy.Name, policy.Service, policy.Path, []byte(`["GET"]`), policy.CreatedAt, policy.UpdatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
		},
		{
			name:        "create error",
			inputPolicy: policy,
			expectError: sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO policies (id, user_id, name, service, path, methods, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?);")).
					WithArgs(policy.ID, policy.UserID, policy.Name, policy.Service, policy.Path, []byte(`["GET"]`), policy.CreatedAt, policy.UpdatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:        "no policy",
			inputPolicy: nil,
			expectError: database.ErrRequiredPolicy,
			setMockDB:   func(mock sqlmock.Sqlmock) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := test.NewMockDB(t)
			defer db.Close()

			ctx := context.Background()

			tt.setMockDB(mock)

			r := database.NewPolicyDBRepository(db)
			if err := r.Create(ctx, tt.inputPolicy); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func TestPolicy_Update(t *testing.T) {
	policyWithoutAgents, err := entity.NewPolicy(uuid.New(), "name", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	agent, err := entity.NewAgent(policyWithoutAgents.UserID, "name")
	if err != nil {
		t.Error(err.Error())
	}
	policyWithAgents, err := entity.NewPolicy(policyWithoutAgents.UserID, "name", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}
	policyWithAgents.SetAgents([]*entity.Agent{agent})

	tests := []struct {
		name        string
		inputPolicy *entity.Policy
		expectError error
		setMockDB   func(sqlmock.Sqlmock)
	}{
		{
			name:        "without permissions",
			inputPolicy: policyWithoutAgents,
			expectError: nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE policies SET user_id = ?, name = ?, service = ?, path = ?, methods = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
					WithArgs(policyWithoutAgents.UserID, policyWithoutAgents.Name, policyWithoutAgents.Service, policyWithoutAgents.Path, []byte(`["GET"]`), policyWithoutAgents.UpdatedAt, policyWithoutAgents.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM permissions WHERE policy_id = ?;")).
					WithArgs(policyWithoutAgents.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
		},
		{
			name:        "with permissions",
			inputPolicy: policyWithAgents,
			expectError: nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE policies SET user_id = ?, name = ?, service = ?, path = ?, methods = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
					WithArgs(policyWithAgents.UserID, policyWithAgents.Name, policyWithAgents.Service, policyWithAgents.Path, []byte(`["GET"]`), policyWithAgents.UpdatedAt, policyWithAgents.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM permissions WHERE policy_id = ?;")).
					WithArgs(policyWithAgents.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO permissions (agent_id, policy_id) VALUES (?, ?);")).
					WithArgs(policyWithAgents.Agents[0], policyWithAgents.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
		},
		{
			name:        "update policy error",
			inputPolicy: policyWithoutAgents,
			expectError: sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE policies SET user_id = ?, name = ?, service = ?, path = ?, methods = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
					WithArgs(policyWithoutAgents.UserID, policyWithoutAgents.Name, policyWithoutAgents.Service, policyWithoutAgents.Path, []byte(`["GET"]`), policyWithoutAgents.UpdatedAt, policyWithoutAgents.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:        "delete permissions error",
			inputPolicy: policyWithoutAgents,
			expectError: sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE policies SET user_id = ?, name = ?, service = ?, path = ?, methods = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
					WithArgs(policyWithoutAgents.UserID, policyWithoutAgents.Name, policyWithoutAgents.Service, policyWithoutAgents.Path, []byte(`["GET"]`), policyWithoutAgents.UpdatedAt, policyWithoutAgents.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM permissions WHERE policy_id = ?;")).
					WithArgs(policyWithoutAgents.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:        "create permissions error",
			inputPolicy: policyWithAgents,
			expectError: sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE policies SET user_id = ?, name = ?, service = ?, path = ?, methods = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
					WithArgs(policyWithAgents.UserID, policyWithAgents.Name, policyWithAgents.Service, policyWithAgents.Path, []byte(`["GET"]`), policyWithAgents.UpdatedAt, policyWithAgents.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM permissions WHERE policy_id = ?;")).
					WithArgs(policyWithAgents.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO permissions (agent_id, policy_id) VALUES (?, ?);")).
					WithArgs(policyWithAgents.Agents[0], policyWithAgents.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:        "no policy",
			inputPolicy: nil,
			expectError: database.ErrRequiredPolicy,
			setMockDB:   func(mock sqlmock.Sqlmock) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := test.NewMockDB(t)
			defer db.Close()

			ctx := context.Background()

			tt.setMockDB(mock)

			r := database.NewPolicyDBRepository(db)
			if err := r.Update(ctx, tt.inputPolicy); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func TestPolicy_Delete(t *testing.T) {
	policy, err := entity.NewPolicy(uuid.New(), "name", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name        string
		inputPolicy *entity.Policy
		expectError error
		setMockDB   func(sqlmock.Sqlmock)
	}{
		{
			name:        "success",
			inputPolicy: policy,
			expectError: nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE policies SET updated_at = updated_at, deleted_at = NOW(6) WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
					WithArgs(policy.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
		},
		{
			name:        "delete error",
			inputPolicy: policy,
			expectError: sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE policies SET updated_at = updated_at, deleted_at = NOW(6) WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
					WithArgs(policy.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:        "no policy",
			inputPolicy: nil,
			expectError: database.ErrRequiredPolicy,
			setMockDB:   func(mock sqlmock.Sqlmock) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := test.NewMockDB(t)
			defer db.Close()

			ctx := context.Background()

			tt.setMockDB(mock)

			r := database.NewPolicyDBRepository(db)
			if err := r.Delete(ctx, tt.inputPolicy); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func TestPolicy_FindOneByIDAndUserIDAndNotDeleted(t *testing.T) {
	policy, err := entity.NewPolicy(uuid.New(), "name", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name         string
		inputID      uuid.UUID
		inputUserID  uuid.UUID
		expectResult *entity.Policy
		expectError  error
		setMockDB    func(sqlmock.Sqlmock)
	}{
		{
			name:         "found",
			inputID:      policy.ID,
			inputUserID:  policy.UserID,
			expectResult: policy,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
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
				)).
					WithArgs(policy.ID, policy.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "service", "path", "methods", "created_at", "updated_at"}).
							AddRow(policy.ID, policy.UserID, policy.Name, policy.Service, policy.Path, fmt.Sprintf(`["%s"]`, strings.Join(policy.Methods, ",")), policy.CreatedAt, policy.UpdatedAt),
					).
					WillReturnError(nil)
			},
		},
		{
			name:         "not found",
			inputID:      policy.ID,
			inputUserID:  policy.UserID,
			expectResult: nil,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
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
				)).
					WithArgs(policy.ID, policy.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "service", "path", "methods", "created_at", "updated_at"}).
							AddRow(policy.ID, policy.UserID, policy.Name, policy.Service, policy.Path, fmt.Sprintf(`["%s"]`, strings.Join(policy.Methods, ",")), policy.CreatedAt, policy.UpdatedAt),
					).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:         "find error",
			inputID:      policy.ID,
			inputUserID:  policy.UserID,
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
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
				)).
					WithArgs(policy.ID, policy.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "service", "path", "methods", "created_at", "updated_at"}).
							AddRow(policy.ID, policy.UserID, policy.Name, policy.Service, policy.Path, fmt.Sprintf(`["%s"]`, strings.Join(policy.Methods, ",")), policy.CreatedAt, policy.UpdatedAt),
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

			r := database.NewPolicyDBRepository(db)
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

func TestPolicy_FindByUserIDAndNotDeleted(t *testing.T) {
	policy, err := entity.NewPolicy(uuid.New(), "name", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name         string
		inputUserID  uuid.UUID
		expectResult []*entity.Policy
		expectError  error
		setMockDB    func(sqlmock.Sqlmock)
	}{
		{
			name:         "found",
			inputUserID:  policy.UserID,
			expectResult: []*entity.Policy{policy},
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, service, path, methods, created_at, updated_at FROM policies WHERE user_id = ? AND deleted_at IS NULL;")).
					WithArgs(policy.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "service", "path", "methods", "created_at", "updated_at"}).
							AddRow(policy.ID, policy.UserID, policy.Name, policy.Service, policy.Path, fmt.Sprintf(`["%s"]`, strings.Join(policy.Methods, ",")), policy.CreatedAt, policy.UpdatedAt),
					).
					WillReturnError(nil)
			},
		},
		{
			name:         "not found",
			inputUserID:  policy.UserID,
			expectResult: []*entity.Policy{},
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, service, path, methods, created_at, updated_at FROM policies WHERE user_id = ? AND deleted_at IS NULL;")).
					WithArgs(policy.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "service", "path", "methods", "created_at", "updated_at"}),
					).
					WillReturnError(nil)
			},
		},
		{
			name:         "find error",
			inputUserID:  policy.UserID,
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, service, path, methods, created_at, updated_at FROM policies WHERE user_id = ? AND deleted_at IS NULL;")).
					WithArgs(policy.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "service", "path", "methods", "created_at", "updated_at"}),
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

			r := database.NewPolicyDBRepository(db)
			result, err := r.FindByUserIDAndNotDeleted(ctx, tt.inputUserID)
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

func TestPolicy_FindByIDsAndUserIDAndNotDeleted(t *testing.T) {
	policy, err := entity.NewPolicy(uuid.New(), "name", "STORAGE", "/", []string{"GET"})
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name         string
		inputIDs     []uuid.UUID
		inputUserID  uuid.UUID
		expectResult []*entity.Policy
		expectError  error
		setMockDB    func(sqlmock.Sqlmock)
	}{
		{
			name:         "found",
			inputIDs:     []uuid.UUID{policy.ID},
			inputUserID:  policy.UserID,
			expectResult: []*entity.Policy{policy},
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, service, path, methods, created_at, updated_at FROM policies WHERE id IN (?) AND user_id = ? AND deleted_at IS NULL;")).
					WithArgs(policy.ID, policy.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "service", "path", "methods", "created_at", "updated_at"}).
							AddRow(policy.ID, policy.UserID, policy.Name, policy.Service, policy.Path, fmt.Sprintf(`["%s"]`, strings.Join(policy.Methods, ",")), policy.CreatedAt, policy.UpdatedAt),
					).
					WillReturnError(nil)
			},
		},
		{
			name:         "not found",
			inputIDs:     []uuid.UUID{policy.ID},
			inputUserID:  policy.UserID,
			expectResult: []*entity.Policy{},
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, service, path, methods, created_at, updated_at FROM policies WHERE id IN (?) AND user_id = ? AND deleted_at IS NULL;")).
					WithArgs(policy.ID, policy.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "service", "path", "methods", "created_at", "updated_at"}),
					).
					WillReturnError(nil)
			},
		},
		{
			name:         "find error",
			inputIDs:     []uuid.UUID{policy.ID},
			inputUserID:  policy.UserID,
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, service, path, methods, created_at, updated_at FROM policies WHERE id IN (?) AND user_id = ? AND deleted_at IS NULL;")).
					WithArgs(policy.ID, policy.UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "name", "service", "path", "methods", "created_at", "updated_at"}),
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

			r := database.NewPolicyDBRepository(db)
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
