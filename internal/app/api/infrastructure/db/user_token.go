package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/domain/repository"
	"holos-auth-api/internal/app/api/infrastructure/model"
	"holos-auth-api/internal/app/api/pkg/apierr"
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

func (i *userTokenInfrastructure) Save(ctx context.Context, userToken *entity.UserToken) apierr.ApiError {
	driver := getSqlxDriver(ctx, i.db)
	userTokenModel := i.convertToModel(userToken)
	if _, err := driver.NamedExecContext(
		ctx,
		`REPLACE user_tokens (user_id, token, expires_at) VALUES (:user_id, :token, :expires_at);`,
		userTokenModel,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (i *userTokenInfrastructure) Delete(ctx context.Context, userToken *entity.UserToken) apierr.ApiError {
	driver := getSqlxDriver(ctx, i.db)
	userTokenModel := i.convertToModel(userToken)
	if _, err := driver.NamedExecContext(
		ctx,
		`DELETE FROM user_tokens WHERE user_id = :user_id;`,
		userTokenModel,
	); err != nil {
		return apierr.NewApiError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (i *userTokenInfrastructure) FindOneByTokenAndNotExpired(ctx context.Context, token string) (*entity.UserToken, apierr.ApiError) {
	var userToken model.UserTokenModel
	driver := getSqlxDriver(ctx, i.db)
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
	return i.convertToEntity(&userToken), nil
}

func (i *userTokenInfrastructure) convertToModel(userToken *entity.UserToken) *model.UserTokenModel {
	return model.NewUserTokenModel(userToken.UserID, userToken.Token, userToken.ExpiresAt)
}

func (i *userTokenInfrastructure) convertToEntity(userToken *model.UserTokenModel) *entity.UserToken {
	return entity.RestoreUserToken(userToken.UserID, userToken.Token, userToken.ExpiresAt)
}
