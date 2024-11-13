package request

type SigninRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
