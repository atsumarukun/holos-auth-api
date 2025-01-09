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

func TestUser_Create(t *testing.T) {
	user, err := entity.NewUser("name", "password", "password")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name        string
		inputUser   *entity.User
		expectError error
		setMockDB   func(sqlmock.Sqlmock)
	}{
		{
			name:        "success",
			inputUser:   user,
			expectError: nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (id, name, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?);")).
					WithArgs(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
		},
		{
			name:        "create error",
			inputUser:   user,
			expectError: sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (id, name, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?);")).
					WithArgs(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:        "no user",
			inputUser:   nil,
			expectError: database.ErrRequiredUser,
			setMockDB:   func(mock sqlmock.Sqlmock) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := test.NewMockDB(t)
			defer db.Close()

			ctx := context.Background()

			tt.setMockDB(mock)

			r := database.NewUserDBRepository(db)
			if err := r.Create(ctx, tt.inputUser); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func TestUser_Update(t *testing.T) {
	user, err := entity.NewUser("name", "password", "password")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name        string
		inputUser   *entity.User
		expectError error
		setMockDB   func(sqlmock.Sqlmock)
	}{
		{
			name:        "success",
			inputUser:   user,
			expectError: nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET name = ?, password = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
					WithArgs(user.Name, user.Password, user.UpdatedAt, user.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
		},
		{
			name:        "update error",
			inputUser:   user,
			expectError: sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET name = ?, password = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
					WithArgs(user.Name, user.Password, user.UpdatedAt, user.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:        "no user",
			inputUser:   nil,
			expectError: database.ErrRequiredUser,
			setMockDB:   func(mock sqlmock.Sqlmock) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := test.NewMockDB(t)
			defer db.Close()

			ctx := context.Background()

			tt.setMockDB(mock)

			r := database.NewUserDBRepository(db)
			if err := r.Update(ctx, tt.inputUser); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func TestUser_Delete(t *testing.T) {
	user, err := entity.NewUser("name", "password", "password")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name        string
		inputUser   *entity.User
		expectError error
		setMockDB   func(sqlmock.Sqlmock)
	}{
		{
			name:        "success",
			inputUser:   user,
			expectError: nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET updated_at = updated_at, deleted_at = NOW(6) WHERE id = ? AND deleted_at IS NULL LIMIT 1;`)).
					WithArgs(user.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(nil)
			},
		},
		{
			name:        "delete error",
			inputUser:   user,
			expectError: sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET updated_at = updated_at, deleted_at = NOW(6) WHERE id = ? AND deleted_at IS NULL LIMIT 1;`)).
					WithArgs(user.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			name:        "no user",
			inputUser:   nil,
			expectError: database.ErrRequiredUser,
			setMockDB:   func(mock sqlmock.Sqlmock) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := test.NewMockDB(t)
			defer db.Close()

			ctx := context.Background()

			tt.setMockDB(mock)

			r := database.NewUserDBRepository(db)
			if err := r.Delete(ctx, tt.inputUser); !errors.Is(err, tt.expectError) {
				t.Errorf("\nexpect: %v\ngot: %v", tt.expectError, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func TestUser_FindOneByIDAndNotDeleted(t *testing.T) {
	user, err := entity.NewUser("name", "password", "password")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name         string
		inputID      uuid.UUID
		expectResult *entity.User
		expectError  error
		setMockDB    func(sqlmock.Sqlmock)
	}{
		{
			name:         "found",
			inputID:      user.ID,
			expectResult: user,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, password, created_at, updated_at FROM users WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
					WithArgs(user.ID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "password", "created_at", "updated_at"}).
							AddRow(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt),
					).
					WillReturnError(nil)
			},
		},
		{
			name:         "not found",
			inputID:      user.ID,
			expectResult: nil,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, password, created_at, updated_at FROM users WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
					WithArgs(user.ID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "password", "created_at", "updated_at"}).
							AddRow(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt),
					).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:         "find error",
			inputID:      user.ID,
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, password, created_at, updated_at FROM users WHERE id = ? AND deleted_at IS NULL LIMIT 1;")).
					WithArgs(user.ID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "password", "created_at", "updated_at"}).
							AddRow(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt),
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

			r := database.NewUserDBRepository(db)
			result, err := r.FindOneByIDAndNotDeleted(ctx, tt.inputID)
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

func TestUser_FindOneByName(t *testing.T) {
	user, err := entity.NewUser("name", "password", "password")
	if err != nil {
		t.Error(err.Error())
	}

	tests := []struct {
		name         string
		inputName    string
		expectResult *entity.User
		expectError  error
		setMockDB    func(sqlmock.Sqlmock)
	}{
		{
			name:         "found",
			inputName:    user.Name,
			expectResult: user,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, password, created_at, updated_at FROM users WHERE name = ? LIMIT 1;")).
					WithArgs(user.Name).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "password", "created_at", "updated_at"}).
							AddRow(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt),
					).
					WillReturnError(nil)
			},
		},
		{
			name:         "not found",
			inputName:    user.Name,
			expectResult: nil,
			expectError:  nil,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, password, created_at, updated_at FROM users WHERE name = ? LIMIT 1;")).
					WithArgs(user.Name).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "password", "created_at", "updated_at"}).
							AddRow(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt),
					).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:         "find error",
			inputName:    user.Name,
			expectResult: nil,
			expectError:  sql.ErrConnDone,
			setMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, password, created_at, updated_at FROM users WHERE name = ? LIMIT 1;")).
					WithArgs(user.Name).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "password", "created_at", "updated_at"}).
							AddRow(user.ID, user.Name, user.Password, user.CreatedAt, user.UpdatedAt),
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

			r := database.NewUserDBRepository(db)
			result, err := r.FindOneByName(ctx, tt.inputName)
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
