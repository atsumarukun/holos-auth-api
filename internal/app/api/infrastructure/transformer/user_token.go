package transformer

import (
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/infrastructure/model"
)

func ToUserTokenModel(userToken *entity.UserToken) *model.UserTokenModel {
	return &model.UserTokenModel{
		UserID:    userToken.UserID,
		Token:     userToken.Token,
		ExpiresAt: userToken.ExpiresAt,
	}
}

func ToUserTokenEntity(userToken *model.UserTokenModel) *entity.UserToken {
	return entity.RestoreUserToken(
		userToken.UserID,
		userToken.Token,
		userToken.ExpiresAt,
	)
}
