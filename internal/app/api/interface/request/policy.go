package request

import "github.com/google/uuid"

type CreatePolicyRequest struct {
	Name    string   `json:"name"`
	Service string   `json:"service"`
	Path    string   `json:"path"`
	Methods []string `json:"methods"`
}

type UpdatePolicyRequest struct {
	Name    string   `json:"name"`
	Service string   `json:"service"`
	Path    string   `json:"path"`
	Methods []string `json:"methods"`
}

type UpdatePolicyAgentsRequest struct {
	AgentIDs []uuid.UUID `json:"agent_ids"`
}
