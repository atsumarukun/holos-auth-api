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

var (
	ErrRequiredUser = status.Error(http.StatusInternalServerError, "user is required")
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
	if user == nil {
		return ErrRequiredUser
	}

	driver := getSqlxDriver(ctx, r.db)
	userModel := transformer.ToUserModel(user)

	_, err := driver.NamedExecContext(
		ctx,
		`INSERT INTO users (id, name, password, created_at, updated_at) VALUES (:id, :name, :password, :created_at, :updated_at);`,
		userModel,
	)

	return err
}

func (r *userDBRepository) Update(ctx context.Context, user *entity.User) error {
	if user == nil {
		return ErrRequiredUser
	}

	driver := getSqlxDriver(ctx, r.db)
	userModel := transformer.ToUserModel(user)

	_, err := driver.NamedExecContext(
		ctx,
		`UPDATE users SET name = :name, password = :password, updated_at = :updated_at WHERE id = :id AND deleted_at IS NULL LIMIT 1;`,
		userModel,
	)

	return err
}

func (r *userDBRepository) Delete(ctx context.Context, user *entity.User) error {
	if user == nil {
		return ErrRequiredUser
	}

	driver := getSqlxDriver(ctx, r.db)
	userModel := transformer.ToUserModel(user)

	_, err := driver.NamedExecContext(
		ctx,
		`UPDATE users SET updated_at = updated_at, deleted_at = NOW(6) WHERE id = :id AND deleted_at IS NULL LIMIT 1;`,
		userModel,
	)

	return err
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
		}
		return nil, err
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
		}
		return nil, err
	}

	return transformer.ToUesrEntity(&user), nil
}
