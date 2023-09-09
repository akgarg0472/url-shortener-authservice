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
	StatusCode int16   `json:"statusCode"`
	UserId     string  `json:"userId"`
	Token      string  `json:"token"`
	Expiration float64 `json:"expiration"`
	Success    bool    `json:"success"`
}
