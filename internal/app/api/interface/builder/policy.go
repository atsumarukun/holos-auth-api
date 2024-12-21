package builder

import (
	"holos-auth-api/internal/app/api/interface/response"
	"holos-auth-api/internal/app/api/usecase/dto"
)

func ToPolicyResponse(policy *dto.PolicyDTO) *response.PolicyResponse {
	return &response.PolicyResponse{
		ID:        policy.ID,
		Name:      policy.Name,
		Service:   policy.Service,
		Path:      policy.Path,
		Methods:   policy.Methods,
		CreatedAt: policy.CreatedAt,
		UpdatedAt: policy.UpdatedAt,
	}
}

func ToPolicyResponses(policies []*dto.PolicyDTO) []*response.PolicyResponse {
	responses := make([]*response.PolicyResponse, len(policies))
	for i, policy := range policies {
		responses[i] = ToPolicyResponse(policy)
	}
	return responses
}
