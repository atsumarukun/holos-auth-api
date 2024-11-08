package infrastructure_test

import (
	"context"
	"database/sql"
	"fmt"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/infrastructure"
	"holos-auth-api/test"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestUser_Create(t *testing.T) {
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
			user, err := entity.NewUser(tt.name, "password", "password")
			if err != nil {
				t.Error(err.Error())
			}

			ctx := context.Background()

			db, mock := test.NewMockDB(t)
			defer db.Close()

			if tt.isTransaction {
				mock.ExpectBegin()
			}
			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (id, name, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?);")).
				WithArgs(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt).
				WillReturnResult(sqlmock.NewResult(1, 1))
			if tt.isTransaction {
				mock.ExpectCommit()
			}

			ui := infrastructure.NewUserInfrastructure(db)
			if tt.isTransaction {
				to := infrastructure.NewSqlxTransactionObject(db)
				if err := to.Transaction(ctx, func(ctx context.Context) error {
					return ui.Create(ctx, user)
				}); err != nil {
					t.Error(err.Error())
				}
			} else {
				if err := ui.Create(ctx, user); err != nil {
					t.Error(err.Error())
				}
			}
		})
	}
}

func TestUser_Update(t *testing.T) {
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
			user, err := entity.NewUser(tt.name, "password", "password")
			if err != nil {
				t.Error(err.Error())
			}

			ctx := context.Background()

			db, mock := test.NewMockDB(t)
			defer db.Close()

			if tt.isTransaction {
				mock.ExpectBegin()
			}
			mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET name = ?, password = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
				WithArgs(user.Name, user.Password, user.UpdatedAt, user.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))
			if tt.isTransaction {
				mock.ExpectCommit()
			}

			ui := infrastructure.NewUserInfrastructure(db)
			if tt.isTransaction {
				to := infrastructure.NewSqlxTransactionObject(db)
				if err := to.Transaction(ctx, func(ctx context.Context) error {
					return ui.Update(ctx, user)
				}); err != nil {
					t.Error(err.Error())
				}
			} else {
				if err := ui.Update(ctx, user); err != nil {
					t.Error(err.Error())
				}
			}
		})
	}
}

func TestUser_Delete(t *testing.T) {
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
			user, err := entity.NewUser(tt.name, "password", "password")
			if err != nil {
				t.Error(err.Error())
			}

			ctx := context.Background()

			db, mock := test.NewMockDB(t)
			defer db.Close()

			if tt.isTransaction {
				mock.ExpectBegin()
			}
			mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET updated_at = ?, deleted_at = NOW(6) WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
				WithArgs(user.UpdatedAt, user.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))
			if tt.isTransaction {
				mock.ExpectCommit()
			}

			ui := infrastructure.NewUserInfrastructure(db)
			if tt.isTransaction {
				to := infrastructure.NewSqlxTransactionObject(db)
				if err := to.Transaction(ctx, func(ctx context.Context) error {
					return ui.Delete(ctx, user)
				}); err != nil {
					t.Error(err.Error())
				}
			} else {
				if err := ui.Delete(ctx, user); err != nil {
					t.Error(err.Error())
				}
			}
		})
	}
}

func TestUser_FindOneByName(t *testing.T) {
	tests := []struct {
		name          string
		isTransaction bool
		resultIsNil   bool
		resultError   error
	}{
		{
			name:          "without_transaction",
			isTransaction: false,
			resultIsNil:   false,
			resultError:   nil,
		},
		{
			name:          "with_transaction",
			isTransaction: true,
			resultIsNil:   false,
			resultError:   nil,
		},
		{
			name:          "user_not_found",
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
			mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, password, created_at, updated_at FROM users WHERE name = ? AND deleted_at IS NULL LIMIT 1;")).
				WithArgs(tt.name).
				WillReturnRows(
					sqlmock.NewRows([]string{"id", "name", "password", "created_at", "updated_at"}).
						AddRow(uuid.New(), tt.name, "password", time.Now(), time.Now()),
				).
				WillReturnError(tt.resultError)
			if tt.isTransaction {
				mock.ExpectCommit()
			}

			ui := infrastructure.NewUserInfrastructure(db)
			if tt.isTransaction {
				to := infrastructure.NewSqlxTransactionObject(db)
				if err := to.Transaction(ctx, func(ctx context.Context) error {
					result, err := ui.FindOneByName(ctx, tt.name)
					if err != nil {
						return err
					}
					if (result == nil) != tt.resultIsNil {
						return fmt.Errorf("expect %t but got %t", (result == nil), tt.resultIsNil)
					}
					return nil
				}); err != nil {
					t.Error(err.Error())
				}
			} else {
				result, err := ui.FindOneByName(ctx, tt.name)
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
