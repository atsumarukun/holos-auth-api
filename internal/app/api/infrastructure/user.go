package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/repository"
	"holos-auth-api/internal/pkg/apierr"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type userInfrastructure struct {
	db *sqlx.DB
}

func NewUserInfrastructure(db *sqlx.DB) repository.UserRepository {
	return &userInfrastructure{
		db: db,
	}
}

func (ui *userInfrastructure) Create(ctx context.Context, user *entity.User) apierr.ApiError {
	driver := getSqlxDriver(ctx, ui.db)
	if _, err := driver.NamedExecContext(
		ctx,
		`INSERT INTO users (id, name, password, created_at, updated_at) VALUES (:id, :name, :password, :created_at, :updated_at);`,
		user,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (ui *userInfrastructure) Update(ctx context.Context, user *entity.User) apierr.ApiError {
	driver := getSqlxDriver(ctx, ui.db)
	if _, err := driver.NamedExecContext(
		ctx,
		`UPDATE users SET name = :name, password = :password, updated_at = :updated_at WHERE id = :id AND deleted_at IS NULL LIMIT 1;`,
		user,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (ui *userInfrastructure) Delete(ctx context.Context, user *entity.User) apierr.ApiError {
	driver := getSqlxDriver(ctx, ui.db)
	if _, err := driver.NamedExecContext(
		ctx,
		`UPDATE users SET updated_at = :updated_at, deleted_at = NOW(6) WHERE id = :id AND deleted_at IS NULL LIMIT 1;`,
		user,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (ui *userInfrastructure) FindOneByName(ctx context.Context, name string) (*entity.User, apierr.ApiError) {
	var user entity.User
	driver := getSqlxDriver(ctx, ui.db)
	if err := driver.QueryRowxContext(
		ctx,
		`SELECT id, name, password, created_at, updated_at FROM users WHERE name = ? AND deleted_at IS NULL LIMIT 1;`,
		name,
	).StructScan(&user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, apierr.NewApiError(http.StatusInternalServerError, err.Error())
		}
	}
	return &user, nil
}
