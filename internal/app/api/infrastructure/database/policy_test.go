package database_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/infrastructure/database"
	"holos-auth-api/internal/app/api/pkg/status"
	"holos-auth-api/test"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestPolicy_Create(t *testing.T) {
	tests := []struct {
		name          string
		isTransaction bool
	}{
		{
			name:          "without_transaction",
			isTransaction: false,
		},
		{
			name:          "with_transaction",
			isTransaction: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy, err := entity.NewPolicy(uuid.New(), "name", "STORAGE", "/", []string{"GET"})
			if err != nil {
				t.Error(err.Error())
			}

			ctx := context.Background()

			db, mock := test.NewMockDB(t)
			defer db.Close()

			if tt.isTransaction {
				mock.ExpectBegin()
			}
			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO policies (id, user_id, name, service, path, methods, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?);")).
				WithArgs(policy.ID, policy.UserID, policy.Name, policy.Service, policy.Path, []byte(`["GET"]`), policy.CreatedAt, policy.UpdatedAt).
				WillReturnResult(sqlmock.NewResult(1, 1))
			if tt.isTransaction {
				mock.ExpectCommit()
			}

			pr := database.NewPolicyDBRepository(db)
			if tt.isTransaction {
				to := database.NewSqlxTransactionObject(db)
				if err := to.Transaction(ctx, func(ctx context.Context) error {
					return pr.Create(ctx, policy)
				}); err != nil {
					t.Error(err.Error())
				}
			} else {
				if err := pr.Create(ctx, policy); err != nil {
					t.Error(err.Error())
				}
			}
		})
	}
}

func TestPolicy_Update(t *testing.T) {
	tests := []struct {
		name          string
		isTransaction bool
	}{
		{
			name:          "without_transaction",
			isTransaction: false,
		},
		{
			name:          "with_transaction",
			isTransaction: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy, err := entity.NewPolicy(uuid.New(), "name", "STORAGE", "/", []string{"GET"})
			if err != nil {
				t.Error(err.Error())
			}

			ctx := context.Background()

			db, mock := test.NewMockDB(t)
			defer db.Close()

			if tt.isTransaction {
				mock.ExpectBegin()
			}
			mock.ExpectExec(regexp.QuoteMeta("UPDATE policies SET user_id = ?, name = ?, service = ?, path = ?, methods = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
				WithArgs(policy.UserID, policy.Name, policy.Service, policy.Path, []byte(`["GET"]`), policy.UpdatedAt, policy.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec(regexp.QuoteMeta("DELETE FROM permissions WHERE policy_id = ?;")).
				WithArgs(policy.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))
			if tt.isTransaction {
				mock.ExpectCommit()
			}

			pr := database.NewPolicyDBRepository(db)
			if tt.isTransaction {
				to := database.NewSqlxTransactionObject(db)
				if err := to.Transaction(ctx, func(ctx context.Context) error {
					return pr.Update(ctx, policy)
				}); err != nil {
					t.Error(err.Error())
				}
			} else {
				if err := pr.Update(ctx, policy); err != nil {
					t.Error(err.Error())
				}
			}
		})
	}
}

func TestPolicy_Delete(t *testing.T) {
	tests := []struct {
		name          string
		isTransaction bool
	}{
		{
			name:          "without_transaction",
			isTransaction: false,
		},
		{
			name:          "with_transaction",
			isTransaction: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy, err := entity.NewPolicy(uuid.New(), "name", "STORAGE", "/", []string{"GET"})
			if err != nil {
				t.Error(err.Error())
			}

			ctx := context.Background()

			db, mock := test.NewMockDB(t)
			defer db.Close()

			if tt.isTransaction {
				mock.ExpectBegin()
			}
			mock.ExpectExec(regexp.QuoteMeta("UPDATE policies SET updated_at = updated_at, deleted_at = NOW(6) WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
				WithArgs(policy.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))
			if tt.isTransaction {
				mock.ExpectCommit()
			}

			pr := database.NewPolicyDBRepository(db)
			if tt.isTransaction {
				to := database.NewSqlxTransactionObject(db)
				if err := to.Transaction(ctx, func(ctx context.Context) error {
					return pr.Delete(ctx, policy)
				}); err != nil {
					t.Error(err.Error())
				}
			} else {
				if err := pr.Delete(ctx, policy); err != nil {
					t.Error(err.Error())
				}
			}
		})
	}
}

func TestPolicy_FindOneByIDAndUserIDAndNotDeleted(t *testing.T) {
	tests := []struct {
		id            uuid.UUID
		userID        uuid.UUID
		name          string
		isTransaction bool
		resultIsNil   bool
		resultError   error
	}{
		{
			id:            uuid.New(),
			userID:        uuid.New(),
			name:          "without_transaction",
			isTransaction: false,
			resultIsNil:   false,
			resultError:   nil,
		},
		{
			id:            uuid.New(),
			userID:        uuid.New(),
			name:          "with_transaction",
			isTransaction: true,
			resultIsNil:   false,
			resultError:   nil,
		},
		{
			id:            uuid.New(),
			userID:        uuid.New(),
			name:          "agent_not_found",
			isTransaction: false,
			resultIsNil:   true,
			resultError:   sql.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			db, mock := test.NewMockDB(t)
			defer db.Close()

			if tt.isTransaction {
				mock.ExpectBegin()
			}
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
				WithArgs(tt.id, tt.userID).
				WillReturnRows(
					sqlmock.NewRows([]string{"id", "user_id", "name", "service", "path", "methods", "created_at", "updated_at"}).
						AddRow(tt.id, tt.userID, tt.name, "STORAGE", "/", []byte(`["GET"]`), time.Now(), time.Now()),
				).
				WillReturnError(tt.resultError)
			if tt.isTransaction {
				mock.ExpectCommit()
			}

			pr := database.NewPolicyDBRepository(db)
			if tt.isTransaction {
				to := database.NewSqlxTransactionObject(db)
				if err := to.Transaction(ctx, func(ctx context.Context) error {
					result, err := pr.FindOneByIDAndUserIDAndNotDeleted(ctx, tt.id, tt.userID)
					if err != nil {
						return err
					}
					if (result == nil) != tt.resultIsNil {
						return status.Error(http.StatusInternalServerError, fmt.Sprintf("expect %t but got %t", (result == nil), tt.resultIsNil))
					}
					return nil
				}); err != nil {
					t.Error(err.Error())
				}
			} else {
				result, err := pr.FindOneByIDAndUserIDAndNotDeleted(ctx, tt.id, tt.userID)
				if err != nil {
					t.Error(err.Error())
				}
				if (result == nil) != tt.resultIsNil {
					t.Errorf("expect %t but got %t", (result == nil), tt.resultIsNil)
				}
			}
		})
	}
}

func TestPolicy_FindByUserIDAndNotDeleted(t *testing.T) {
	tests := []struct {
		id            uuid.UUID
		userID        uuid.UUID
		name          string
		isTransaction bool
		resultIsNil   bool
		resultError   error
	}{
		{
			id:            uuid.New(),
			userID:        uuid.New(),
			name:          "without_transaction",
			isTransaction: false,
			resultIsNil:   false,
			resultError:   nil,
		},
		{
			id:            uuid.New(),
			userID:        uuid.New(),
			name:          "with_transaction",
			isTransaction: true,
			resultIsNil:   false,
			resultError:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			db, mock := test.NewMockDB(t)
			defer db.Close()

			if tt.isTransaction {
				mock.ExpectBegin()
			}
			mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, service, path, methods, created_at, updated_at FROM policies WHERE user_id = ? AND deleted_at IS NULL;")).
				WithArgs(tt.userID).
				WillReturnRows(
					sqlmock.NewRows([]string{"id", "user_id", "name", "service", "path", "methods", "created_at", "updated_at"}).
						AddRow(tt.id, tt.userID, tt.name, "STORAGE", "/", []byte(`["GET"]`), time.Now(), time.Now()),
				).
				WillReturnError(tt.resultError)
			if tt.isTransaction {
				mock.ExpectCommit()
			}

			pr := database.NewPolicyDBRepository(db)
			if tt.isTransaction {
				to := database.NewSqlxTransactionObject(db)
				if err := to.Transaction(ctx, func(ctx context.Context) error {
					result, err := pr.FindByUserIDAndNotDeleted(ctx, tt.userID)
					if err != nil {
						return err
					}
					if (result == nil) != tt.resultIsNil {
						return status.Error(http.StatusInternalServerError, fmt.Sprintf("expect %t but got %t", (result == nil), tt.resultIsNil))
					}
					return nil
				}); err != nil {
					t.Error(err.Error())
				}
			} else {
				result, err := pr.FindByUserIDAndNotDeleted(ctx, tt.userID)
				if err != nil {
					t.Error(err.Error())
				}
				if (result == nil) != tt.resultIsNil {
					t.Errorf("expect %t but got %t", (result == nil), tt.resultIsNil)
				}
			}
		})
	}
}

func TestPolicy_FindByIDsAndUserIDAndNotDeleted(t *testing.T) {
	tests := []struct {
		id            uuid.UUID
		ids           []uuid.UUID
		userID        uuid.UUID
		name          string
		isTransaction bool
		resultIsNil   bool
		resultError   error
	}{
		{
			id:            uuid.New(),
			ids:           []uuid.UUID{uuid.New()},
			userID:        uuid.New(),
			name:          "without_transaction",
			isTransaction: false,
			resultIsNil:   false,
			resultError:   nil,
		},
		{
			id:            uuid.New(),
			ids:           []uuid.UUID{uuid.New()},
			userID:        uuid.New(),
			name:          "with_transaction",
			isTransaction: true,
			resultIsNil:   false,
			resultError:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			db, mock := test.NewMockDB(t)
			defer db.Close()

			args := make([]driver.Value, len(tt.ids)+1)
			for i, id := range tt.ids {
				args[i] = id
			}
			args[len(args)-1] = tt.userID

			if tt.isTransaction {
				mock.ExpectBegin()
			}
			mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, service, path, methods, created_at, updated_at FROM policies WHERE id IN (?) AND user_id = ? AND deleted_at IS NULL;")).
				WithArgs(args...).
				WillReturnRows(
					sqlmock.NewRows([]string{"id", "user_id", "name", "service", "path", "methods", "created_at", "updated_at"}).
						AddRow(tt.id, tt.userID, tt.name, "STORAGE", "/", []byte(`["GET"]`), time.Now(), time.Now()),
				).
				WillReturnError(tt.resultError)
			if tt.isTransaction {
				mock.ExpectCommit()
			}

			pr := database.NewPolicyDBRepository(db)
			if tt.isTransaction {
				to := database.NewSqlxTransactionObject(db)
				if err := to.Transaction(ctx, func(ctx context.Context) error {
					result, err := pr.FindByIDsAndUserIDAndNotDeleted(ctx, tt.ids, tt.userID)
					if err != nil {
						return err
					}
					if (result == nil) != tt.resultIsNil {
						return status.Error(http.StatusInternalServerError, fmt.Sprintf("expect %t but got %t", (result == nil), tt.resultIsNil))
					}
					return nil
				}); err != nil {
					t.Error(err.Error())
				}
			} else {
				result, err := pr.FindByIDsAndUserIDAndNotDeleted(ctx, tt.ids, tt.userID)
				if err != nil {
					t.Error(err.Error())
				}
				if (result == nil) != tt.resultIsNil {
					t.Errorf("expect %t but got %t", (result == nil), tt.resultIsNil)
				}
			}
		})
	}
}
