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

var (
	ErrRequiredUserToken = status.Error(http.StatusInternalServerError, "user token is required")
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
	if userToken == nil {
		return ErrRequiredUserToken
	}

	driver := getSqlxDriver(ctx, r.db)
	userTokenModel := transformer.ToUserTokenModel(userToken)

	_, err := driver.NamedExecContext(
		ctx,
		`REPLACE user_tokens (user_id, token, expires_at) VALUES (:user_id, :token, :expires_at);`,
		userTokenModel,
	)

	return err
}

func (r *userTokenDBRepository) Delete(ctx context.Context, userToken *entity.UserToken) error {
	if userToken == nil {
		return ErrRequiredUserToken
	}

	driver := getSqlxDriver(ctx, r.db)
	userTokenModel := transformer.ToUserTokenModel(userToken)

	_, err := driver.NamedExecContext(
		ctx,
		`DELETE FROM user_tokens WHERE user_id = :user_id;`,
		userTokenModel,
	)

	return err
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
		}
		return nil, err
	}

	return transformer.ToUserTokenEntity(&userToken), nil
}
