package request

import "github.com/google/uuid"

type CreateAgentRequest struct {
	Name string `json:"name"`
}

type UpdateAgentRequest struct {
	Name string `json:"name"`
}

type UpdateAgentPoliciesRequest struct {
	PolicyIDs []uuid.UUID `json:"policy_ids"`
}
