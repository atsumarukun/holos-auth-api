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

func TestAgent_Create(t *testing.T) {
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
			agent, err := entity.NewAgent(uuid.New(), "name")
			if err != nil {
				t.Error(err.Error())
			}

			ctx := context.Background()

			db, mock := test.NewMockDB(t)
			defer db.Close()

			if tt.isTransaction {
				mock.ExpectBegin()
			}
			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO agents (id, user_id, name, created_at, updated_at) VALUES (?, ?, ?, ?, ?);")).
				WithArgs(agent.ID, agent.UserID, agent.Name, agent.CreatedAt, agent.UpdatedAt).
				WillReturnResult(sqlmock.NewResult(1, 1))
			if tt.isTransaction {
				mock.ExpectCommit()
			}

			ar := database.NewAgentDBRepository(db)
			if tt.isTransaction {
				to := database.NewSqlxTransactionObject(db)
				if err := to.Transaction(ctx, func(ctx context.Context) error {
					return ar.Create(ctx, agent)
				}); err != nil {
					t.Error(err.Error())
				}
			} else {
				if err := ar.Create(ctx, agent); err != nil {
					t.Error(err.Error())
				}
			}
		})
	}
}

func TestAgent_Update(t *testing.T) {
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
			agent, err := entity.NewAgent(uuid.New(), "name")
			if err != nil {
				t.Error(err.Error())
			}

			ctx := context.Background()

			db, mock := test.NewMockDB(t)
			defer db.Close()

			if tt.isTransaction {
				mock.ExpectBegin()
			}
			mock.ExpectExec(regexp.QuoteMeta("UPDATE agents SET user_id = ?, name = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
				WithArgs(agent.UserID, agent.Name, agent.UpdatedAt, agent.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec(regexp.QuoteMeta("DELETE FROM permissions WHERE agent_id = ?;")).
				WithArgs(agent.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))
			if tt.isTransaction {
				mock.ExpectCommit()
			}

			ar := database.NewAgentDBRepository(db)
			if tt.isTransaction {
				to := database.NewSqlxTransactionObject(db)
				if err := to.Transaction(ctx, func(ctx context.Context) error {
					return ar.Update(ctx, agent)
				}); err != nil {
					t.Error(err.Error())
				}
			} else {
				if err := ar.Update(ctx, agent); err != nil {
					t.Error(err.Error())
				}
			}
		})
	}
}

func TestAgent_Delete(t *testing.T) {
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
			agent, err := entity.NewAgent(uuid.New(), "name")
			if err != nil {
				t.Error(err.Error())
			}

			ctx := context.Background()

			db, mock := test.NewMockDB(t)
			defer db.Close()

			if tt.isTransaction {
				mock.ExpectBegin()
			}
			mock.ExpectExec(regexp.QuoteMeta("UPDATE agents SET updated_at = updated_at, deleted_at = NOW(6) WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
				WithArgs(agent.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))
			if tt.isTransaction {
				mock.ExpectCommit()
			}

			ar := database.NewAgentDBRepository(db)
			if tt.isTransaction {
				to := database.NewSqlxTransactionObject(db)
				if err := to.Transaction(ctx, func(ctx context.Context) error {
					return ar.Delete(ctx, agent)
				}); err != nil {
					t.Error(err.Error())
				}
			} else {
				if err := ar.Delete(ctx, agent); err != nil {
					t.Error(err.Error())
				}
			}
		})
	}
}

func TestAgent_FindOneByIDAndUserIDAndNotDeleted(t *testing.T) {
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
				WithArgs(tt.id, tt.userID).
				WillReturnRows(
					sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}).
						AddRow(tt.id, tt.userID, tt.name, time.Now(), time.Now()),
				).
				WillReturnError(tt.resultError)
			if tt.isTransaction {
				mock.ExpectCommit()
			}

			ar := database.NewAgentDBRepository(db)
			if tt.isTransaction {
				to := database.NewSqlxTransactionObject(db)
				if err := to.Transaction(ctx, func(ctx context.Context) error {
					result, err := ar.FindOneByIDAndUserIDAndNotDeleted(ctx, tt.id, tt.userID)
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
				result, err := ar.FindOneByIDAndUserIDAndNotDeleted(ctx, tt.id, tt.userID)
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

func TestAgent_FindByUserIDAndNotDeleted(t *testing.T) {
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
			mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, created_at, updated_at FROM agents WHERE user_id = ? AND deleted_at IS NULL;")).
				WithArgs(tt.userID).
				WillReturnRows(
					sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}).
						AddRow(tt.id, tt.userID, tt.name, time.Now(), time.Now()),
				).
				WillReturnError(tt.resultError)
			if tt.isTransaction {
				mock.ExpectCommit()
			}

			ar := database.NewAgentDBRepository(db)
			if tt.isTransaction {
				to := database.NewSqlxTransactionObject(db)
				if err := to.Transaction(ctx, func(ctx context.Context) error {
					result, err := ar.FindByUserIDAndNotDeleted(ctx, tt.userID)
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
				result, err := ar.FindByUserIDAndNotDeleted(ctx, tt.userID)
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

func TestAgent_FindByIDsAndUserIDAndNotDeleted(t *testing.T) {
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
			mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, created_at, updated_at FROM agents WHERE id IN (?) AND user_id = ? AND deleted_at IS NULL;")).
				WithArgs(args...).
				WillReturnRows(
					sqlmock.NewRows([]string{"id", "user_id", "name", "created_at", "updated_at"}).
						AddRow(tt.id, tt.userID, tt.name, time.Now(), time.Now()),
				).
				WillReturnError(tt.resultError)
			if tt.isTransaction {
				mock.ExpectCommit()
			}

			ar := database.NewAgentDBRepository(db)
			if tt.isTransaction {
				to := database.NewSqlxTransactionObject(db)
				if err := to.Transaction(ctx, func(ctx context.Context) error {
					result, err := ar.FindByIDsAndUserIDAndNotDeleted(ctx, tt.ids, tt.userID)
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
				result, err := ar.FindByIDsAndUserIDAndNotDeleted(ctx, tt.ids, tt.userID)
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
