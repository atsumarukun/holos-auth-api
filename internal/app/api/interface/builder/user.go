package builder

import (
	"holos-auth-api/internal/app/api/interface/response"
	"holos-auth-api/internal/app/api/usecase/dto"
)

func ToUserResponse(user *dto.UserDTO) *response.UserResponse {
	return &response.UserResponse{
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
