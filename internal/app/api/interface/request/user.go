package request

type CreateUserRequest struct {
	Name            string `json:"name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type UpdateUserNameRequest struct {
	Name string `json:"name"`
}

type UpdateUserPasswordRequest struct {
	CurrentPassword    string `json:"current_password"`
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password"`
}

type DeleteUserRequest struct {
	Password string `json:"password"`
}
