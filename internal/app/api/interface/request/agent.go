package request

type CreateAgentRequest struct {
	Name string `json:"name"`
}

type UpdateAgentRequest struct {
	Name string `json:"name"`
}
