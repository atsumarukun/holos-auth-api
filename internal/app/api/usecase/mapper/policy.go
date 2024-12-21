package mapper

import (
	"holos-auth-api/internal/app/api/domain/entity"
	"holos-auth-api/internal/app/api/usecase/dto"
)

func ToPolicyDTO(policy *entity.Policy) *dto.PolicyDTO {
	return &dto.PolicyDTO{
		ID:        policy.ID,
		UserID:    policy.UserID,
		Name:      policy.Name,
		Service:   policy.Service,
		Path:      policy.Path,
		Methods:   policy.Methods,
		Agents:    policy.Agents,
		CreatedAt: policy.CreatedAt,
		UpdatedAt: policy.UpdatedAt,
	}
}

func ToPolicyDTOs(policies []*entity.Policy) []*dto.PolicyDTO {
	dtos := make([]*dto.PolicyDTO, len(policies))
	for i, policy := range policies {
		dtos[i] = ToPolicyDTO(policy)
	}
	return dtos
}
