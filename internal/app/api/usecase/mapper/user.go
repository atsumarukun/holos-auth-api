package mapper

import (
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/usecase/dto"
)

func ToUserDTO(user *entity.User) *dto.UserDTO {
	return &dto.UserDTO{
		ID:        user.ID,
		Name:      user.Name,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
