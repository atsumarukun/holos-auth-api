package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/repository"

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

func (ui *userInfrastructure) Create(ctx context.Context, user *entity.User) error {
	driver := getSqlxDriver(ctx, ui.db)
	_, err := driver.NamedExecContext(
		ctx,
		`INSERT INTO users (id, name, password, created_at, updated_at) VALUES (:id, :name, :password, :created_at, :updated_at);`,
		user,
	)
	return err
}

func (ui *userInfrastructure) Update(ctx context.Context, user *entity.User) error {
	driver := getSqlxDriver(ctx, ui.db)
	_, err := driver.NamedExecContext(
		ctx,
		`UPDATE users SET name = :name, password = :password, updated_at = :updated_at WHERE id = :id AND deleted_at IS NULL LIMIT 1;`,
		user,
	)
	return err
}

func (ui *userInfrastructure) Delete(ctx context.Context, user *entity.User) error {
	driver := getSqlxDriver(ctx, ui.db)
	_, err := driver.NamedExecContext(
		ctx,
		`UPDATE users SET updated_at = :updated_at, deleted_at = NOW(6) WHERE id = :id AND deleted_at IS NULL LIMIT 1;`,
		user,
	)
	return err
}

func (ui *userInfrastructure) FindOneByName(ctx context.Context, name string) (*entity.User, error) {
	var user entity.User
	driver := getSqlxDriver(ctx, ui.db)
	err := driver.QueryRowxContext(
		ctx,
		`SELECT id, name, password, created_at, updated_at FROM users WHERE name = ? AND deleted_at IS NULL LIMIT 1;`,
		name,
	).StructScan(&user)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &user, err
}
