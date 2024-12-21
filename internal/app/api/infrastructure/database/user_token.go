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

	"github.com/jmoiron/sqlx"
)

type userTokenDBRepository struct {
	db *sqlx.DB
}

func NewUserTokenDBRepository(db *sqlx.DB) repository.UserTokenRepository {
	return &userTokenDBRepository{
		db: db,
	}
}

func (r *userTokenDBRepository) Save(ctx context.Context, userToken *entity.UserToken) error {
	driver := getSqlxDriver(ctx, r.db)
	userTokenModel := transformer.ToUserTokenModel(userToken)
	if _, err := driver.NamedExecContext(
		ctx,
		`REPLACE user_tokens (user_id, token, expires_at) VALUES (:user_id, :token, :expires_at);`,
		userTokenModel,
	); err != nil {
		return status.Error(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *userTokenDBRepository) Delete(ctx context.Context, userToken *entity.UserToken) error {
	driver := getSqlxDriver(ctx, r.db)
	userTokenModel := transformer.ToUserTokenModel(userToken)
	if _, err := driver.NamedExecContext(
		ctx,
		`DELETE FROM user_tokens WHERE user_id = :user_id;`,
		userTokenModel,
	); err != nil {
		return status.Error(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *userTokenDBRepository) FindOneByTokenAndNotExpired(ctx context.Context, token string) (*entity.UserToken, error) {
	var userToken model.UserTokenModel
	driver := getSqlxDriver(ctx, r.db)
	if err := driver.QueryRowxContext(
		ctx,
		`SELECT user_id, token, expires_at FROM user_tokens WHERE token = ? AND NOW(6) < expires_at LIMIT 1;`,
		token,
	).StructScan(&userToken); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}
	}
	return transformer.ToUserTokenEntity(&userToken), nil
}
