package model

type LoginResponse struct {
	AccessToken string `json:"auth_token"`
	UserId      string `json:"user_id"`
	Name        string `json:"name"`
}

type SignupResponse struct {
	Message    string `json:"message"`
	StatusCode int16  `json:"status_code"`
}

type ErrorResponse struct {
	Message   interface{} `json:"message"`
	ErrorCode int16       `json:"error_code"`
	Errors    interface{} `json:"errors"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

type ValidateTokenResponse struct {
	UserId     string  `json:"userId"`
	Token      string  `json:"token"`
	Expiration float64 `json:"expiration"`
	Success    bool    `json:"success"`
}

type ForgotPasswordResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}
