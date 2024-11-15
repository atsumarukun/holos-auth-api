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

type userTokenInfrastructure struct {
	db *sqlx.DB
}

func NewUserTokenInfrastructure(db *sqlx.DB) repository.UserTokenRepository {
	return &userTokenInfrastructure{
		db: db,
	}
}

func (uti *userTokenInfrastructure) Save(ctx context.Context, userToken *entity.UserToken) apierr.ApiError {
	driver := getSqlxDriver(ctx, uti.db)
	if _, err := driver.NamedExecContext(
		ctx,
		`REPLACE user_tokens (user_id, token, expires_at) VALUES (:user_id, :token, :expires_at);`,
		userToken,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (uti *userTokenInfrastructure) Delete(ctx context.Context, userToken *entity.UserToken) apierr.ApiError {
	driver := getSqlxDriver(ctx, uti.db)
	if _, err := driver.NamedExecContext(
		ctx,
		`DELETE FROM user_tokens WHERE user_id = :user_id;`,
		userToken,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (uti *userTokenInfrastructure) FindOneByTokenAndNotExpired(ctx context.Context, token string) (*entity.UserToken, apierr.ApiError) {
	var userToken entity.UserToken
	driver := getSqlxDriver(ctx, uti.db)
	if err := driver.QueryRowxContext(
		ctx,
		`SELECT user_id, token, expires_at FROM user_tokens WHERE token = ? AND NOW(6) < expires_at LIMIT 1;`,
		token,
	).StructScan(&userToken); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, apierr.NewApiError(http.StatusInternalServerError, err.Error())
		}
	}
	return &userToken, nil
}
