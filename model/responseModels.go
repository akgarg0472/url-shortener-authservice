package model

type LoginResponse struct {
	AccessToken string `json:"auth_token"`
}

type SignupResponse struct {
	Message string `json:"message"`
	UserId  string `json:"user_id"`
}

type ErrorResponse struct {
	Message   interface{} `json:"message"`
	ErrorCode int         `json:"error_code"`
}
