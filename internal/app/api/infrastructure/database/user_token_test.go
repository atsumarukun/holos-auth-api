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

func TestUserToken_Save(t *testing.T) {
	userToken, err := entity.NewUserToken(uuid.New())
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name           string
		inputUserToken *entity.UserToken
		expectError    error
		setMockDB      func(sqlmock.Sqlmock)
	}{
		{
			name:           "success",
			inputUserToken: userToken,
			expectError:    nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("REPLACE user_tokens (user_id, token, expires_at) VALUES (?, ?, ?);")).
					WithArgs(userToken.UserID, userToken.Token, userToken.ExpiresAt).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
		},
		{
			name:           "save error",
			inputUserToken: userToken,
			expectError:    sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("REPLACE user_tokens (user_id, token, expires_at) VALUES (?, ?, ?);")).
					WithArgs(userToken.UserID, userToken.Token, userToken.ExpiresAt).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:           "no user token",
			inputUserToken: nil,
			expectError:    database.ErrRequiredUserToken,
			setMockDB:      func(mock sqlmock.Sqlmock) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := test.NewMockDB(t)
			defer db.Close()

			ctx := context.Background()

			tt.setMockDB(mock)

			r := database.NewUserTokenDBRepository(db)
			if err := r.Save(ctx, tt.inputUserToken); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func TestUserToken_Delete(t *testing.T) {
	userToken, err := entity.NewUserToken(uuid.New())
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name           string
		inputUserToken *entity.UserToken
		expectError    error
		setMockDB      func(sqlmock.Sqlmock)
	}{
		{
			name:           "success",
			inputUserToken: userToken,
			expectError:    nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM user_tokens WHERE user_id = ?;")).
					WithArgs(userToken.UserID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
		},
		{
			name:           "delete error",
			inputUserToken: userToken,
			expectError:    sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM user_tokens WHERE user_id = ?;")).
					WithArgs(userToken.UserID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:           "no user token",
			inputUserToken: nil,
			expectError:    database.ErrRequiredUserToken,
			setMockDB:      func(mock sqlmock.Sqlmock) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := test.NewMockDB(t)
			defer db.Close()

			ctx := context.Background()

			tt.setMockDB(mock)

			r := database.NewUserTokenDBRepository(db)
			if err := r.Delete(ctx, tt.inputUserToken); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func TestUserToken_FindOneByTokenAndNotExpired(t *testing.T) {
	userToken, err := entity.NewUserToken(uuid.New())
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name         string
		inputToken   string
		expectResult *entity.UserToken
		expectError  error
		setMockDB    func(sqlmock.Sqlmock)
	}{
		{
			name:         "found",
			inputToken:   userToken.Token,
			expectResult: userToken,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, token, expires_at FROM user_tokens WHERE token = ? AND NOW(6) < expires_at LIMIT 1;")).
					WithArgs(userToken.Token).
					WillReturnRows(
						sqlmock.NewRows([]string{"user_id", "token", "expires_at"}).
							AddRow(userToken.UserID, userToken.Token, userToken.ExpiresAt),
					).
					WillReturnError(nil)
			},
		},
		{
			name:         "not found",
			inputToken:   userToken.Token,
			expectResult: nil,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, token, expires_at FROM user_tokens WHERE token = ? AND NOW(6) < expires_at LIMIT 1;")).
					WithArgs(userToken.Token).
					WillReturnRows(
						sqlmock.NewRows([]string{"user_id", "token", "expires_at"}).
							AddRow(userToken.UserID, userToken.Token, userToken.ExpiresAt),
					).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:         "not found",
			inputToken:   userToken.Token,
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, token, expires_at FROM user_tokens WHERE token = ? AND NOW(6) < expires_at LIMIT 1;")).
					WithArgs(userToken.Token).
					WillReturnRows(
						sqlmock.NewRows([]string{"user_id", "token", "expires_at"}).
							AddRow(userToken.UserID, userToken.Token, userToken.ExpiresAt),
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

			r := database.NewUserTokenDBRepository(db)
			result, err := r.FindOneByTokenAndNotExpired(ctx, tt.inputToken)
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
