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
	ErrorCode int16       `json:"error_code"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

type ValidateTokenResponse struct {
	Message string `json:"message"`
}
