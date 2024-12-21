package transformer

import (
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/infrastructure/model"
)

func ToUserModel(user *entity.User) *model.UserModel {
	return &model.UserModel{
		ID:        user.ID,
		Name:      user.Name,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUesrEntity(user *model.UserModel) *entity.User {
	return entity.RestoreUser(
		user.ID,
		user.Name,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)
}
