package response

import "time"

type UserResponse struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUserResponse(name string, createdAt time.Time, updatedAt time.Time) *UserResponse {
	return &UserResponse{
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
