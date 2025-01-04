package response

import (
	"time"

	"github.com/google/uuid"
)

type PolicyResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Effect    string    `json:"effect"`
	Service   string    `json:"service"`
	Path      string    `json:"path"`
	Methods   []string  `json:"methods"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
