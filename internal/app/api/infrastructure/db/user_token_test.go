package infrastructure_test

import (
	"context"
	"database/sql"
	"fmt"
	"holos-auth-api/internal/app/api/domain/entity"
	dbRepository "holos-auth-api/internal/app/api/infrastructure/db"
	"holos-auth-api/internal/app/api/pkg/apierr"
	"holos-auth-api/test"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestUserToken_Save(t *testing.T) {
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
			userToken, err := entity.NewUserToken(uuid.New())
			if err != nil {
				t.Error(err.Error())
			}

			ctx := context.Background()

			db, mock := test.NewMockDB(t)
			defer db.Close()

			if tt.isTransaction {
				mock.ExpectBegin()
			}
			mock.ExpectExec(regexp.QuoteMeta("REPLACE user_tokens (user_id, token, expires_at) VALUES (?, ?, ?);")).
				WithArgs(userToken.UserID, userToken.Token, userToken.ExpiresAt).
				WillReturnResult(sqlmock.NewResult(1, 1))
			if tt.isTransaction {
				mock.ExpectCommit()
			}

			ur := dbRepository.NewUserTokenDBRepository(db)
			if tt.isTransaction {
				to := dbRepository.NewSqlxTransactionObject(db)
				if err := to.Transaction(ctx, func(ctx context.Context) apierr.ApiError {
					return ur.Save(ctx, userToken)
				}); err != nil {
					t.Error(err.Error())
				}
			} else {
				if err := ur.Save(ctx, userToken); err != nil {
					t.Error(err.Error())
				}
			}
		})
	}
}

func TestUserToken_Delete(t *testing.T) {
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
			userToken, err := entity.NewUserToken(uuid.New())
			if err != nil {
				t.Error(err.Error())
			}

			ctx := context.Background()

			db, mock := test.NewMockDB(t)
			defer db.Close()

			if tt.isTransaction {
				mock.ExpectBegin()
			}
			mock.ExpectExec(regexp.QuoteMeta("DELETE FROM user_tokens WHERE user_id = ?;")).
				WithArgs(userToken.UserID).
				WillReturnResult(sqlmock.NewResult(1, 1))
			if tt.isTransaction {
				mock.ExpectCommit()
			}

			ur := dbRepository.NewUserTokenDBRepository(db)
			if tt.isTransaction {
				to := dbRepository.NewSqlxTransactionObject(db)
				if err := to.Transaction(ctx, func(ctx context.Context) apierr.ApiError {
					return ur.Delete(ctx, userToken)
				}); err != nil {
					t.Error(err.Error())
				}
			} else {
				if err := ur.Delete(ctx, userToken); err != nil {
					t.Error(err.Error())
				}
			}
		})
	}
}

func TestUserToken_FindOneByTokenAndNotExpired(t *testing.T) {
	tests := []struct {
		token         string
		isTransaction bool
		resultIsNil   bool
		resultError   error
	}{
		{
			token:         "without_transaction",
			isTransaction: false,
			resultIsNil:   false,
			resultError:   nil,
		},
		{
			token:         "with_transaction",
			isTransaction: true,
			resultIsNil:   false,
			resultError:   nil,
		},
		{
			token:         "user_not_found",
			isTransaction: false,
			resultIsNil:   true,
			resultError:   sql.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			ctx := context.Background()

			db, mock := test.NewMockDB(t)
			defer db.Close()

			if tt.isTransaction {
				mock.ExpectBegin()
			}
			mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, token, expires_at FROM user_tokens WHERE token = ? AND NOW(6) < expires_at LIMIT 1;")).
				WithArgs(tt.token).
				WillReturnRows(
					sqlmock.NewRows([]string{"user_id", "token", "expires_at"}).
						AddRow(uuid.New(), tt.token, time.Now()),
				).
				WillReturnError(tt.resultError)
			if tt.isTransaction {
				mock.ExpectCommit()
			}

			ur := dbRepository.NewUserTokenDBRepository(db)
			if tt.isTransaction {
				to := dbRepository.NewSqlxTransactionObject(db)
				if err := to.Transaction(ctx, func(ctx context.Context) apierr.ApiError {
					result, err := ur.FindOneByTokenAndNotExpired(ctx, tt.token)
					if err != nil {
						return err
					}
					if (result == nil) != tt.resultIsNil {
						return apierr.NewApiError(http.StatusInternalServerError, fmt.Sprintf("expect %t but got %t", (result == nil), tt.resultIsNil))
					}
					return nil
				}); err != nil {
					t.Error(err.Error())
				}
			} else {
				result, err := ur.FindOneByTokenAndNotExpired(ctx, tt.token)
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
