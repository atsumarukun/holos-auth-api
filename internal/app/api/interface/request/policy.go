package request

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
