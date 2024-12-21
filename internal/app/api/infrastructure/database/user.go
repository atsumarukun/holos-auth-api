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

type userDBRepository struct {
	db *sqlx.DB
}

func NewUserDBRepository(db *sqlx.DB) repository.UserRepository {
	return &userDBRepository{
		db: db,
	}
}

func (r *userDBRepository) Create(ctx context.Context, user *entity.User) error {
	driver := getSqlxDriver(ctx, r.db)
	userModel := transformer.ToUserModel(user)
	if _, err := driver.NamedExecContext(
		ctx,
		`INSERT INTO users (id, name, password, created_at, updated_at) VALUES (:id, :name, :password, :created_at, :updated_at);`,
		userModel,
	); err != nil {
		return status.Error(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *userDBRepository) Update(ctx context.Context, user *entity.User) error {
	driver := getSqlxDriver(ctx, r.db)
	userModel := transformer.ToUserModel(user)
	if _, err := driver.NamedExecContext(
		ctx,
		`UPDATE users SET name = :name, password = :password, updated_at = :updated_at WHERE id = :id AND deleted_at IS NULL LIMIT 1;`,
		userModel,
	); err != nil {
		return status.Error(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *userDBRepository) Delete(ctx context.Context, user *entity.User) error {
	driver := getSqlxDriver(ctx, r.db)
	userModel := transformer.ToUserModel(user)
	if _, err := driver.NamedExecContext(
		ctx,
		`UPDATE users
		LEFT JOIN agents ON users.id = agents.user_id
		SET
			users.updated_at = users.updated_at,
			users.deleted_at = NOW(6),
			agents.updated_at = agents.updated_at,
			agents.deleted_at = NOW(6)
		WHERE
			users.id = :id
			AND users.deleted_at IS NULL
			AND agents.deleted_at IS NULL;`,
		userModel,
	); err != nil {
		return status.Error(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *userDBRepository) FindOneByIDAndNotDeleted(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user model.UserModel
	driver := getSqlxDriver(ctx, r.db)
	if err := driver.QueryRowxContext(
		ctx,
		`SELECT id, name, password, created_at, updated_at FROM users WHERE id = ? AND deleted_at IS NULL LIMIT 1;`,
		id,
	).StructScan(&user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	}
	return transformer.ToUesrEntity(&user), nil
}

func (r *userDBRepository) FindOneByName(ctx context.Context, name string) (*entity.User, error) {
	var user model.UserModel
	driver := getSqlxDriver(ctx, r.db)
	if err := driver.QueryRowxContext(
		ctx,
		`SELECT id, name, password, created_at, updated_at FROM users WHERE name = ? LIMIT 1;`,
		name,
	).StructScan(&user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	}
	return transformer.ToUesrEntity(&user), nil
}
