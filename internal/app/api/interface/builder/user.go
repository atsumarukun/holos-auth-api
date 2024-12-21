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

func ToUserResponses(users []*dto.UserDTO) []*response.UserResponse {
	responses := make([]*response.UserResponse, len(users))
	for i, user := range users {
		responses[i] = ToUserResponse(user)
	}
	return responses
}
