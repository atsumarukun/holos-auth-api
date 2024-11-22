package request

type CreatePolicyRequest struct {
	Name           string   `json:"name"`
	Service        string   `json:"service"`
	Path           string   `json:"path"`
	AllowedMethods []string `json:"allowed_methods"`
}

type UpdatePolicyRequest struct {
	Name           string   `json:"name"`
	Service        string   `json:"service"`
	Path           string   `json:"path"`
	AllowedMethods []string `json:"allowed_methods"`
}
