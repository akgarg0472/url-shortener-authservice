package model

import "fmt"

type LoginResponse struct {
	AccessToken string `json:"auth_token"`
	UserId      string `json:"user_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
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

type ResetPasswordResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type OAuthClient struct {
	Provider    string `json:"provider"`
	ClientId    string `json:"client_id"`
	RedirectURI string `json:"redirect_uri"`
	AccessType  string `json:"access_type"`
	Scope       string `json:"scope"`
}

type OAuthClientResponse struct {
	Clients    []OAuthClient `json:"clients"`
	Success    bool          `json:"success"`
	StatusCode int           `json:"status_code"`
}

type OAuthCallbackResponse struct {
	Success   bool   `json:"success"`
	UserId    string `json:"user_id"`
	AuthToken string `json:"auth_token"`
}

func (c OAuthClient) String() string {
	return fmt.Sprintf("OAuthClient{Provider: %s, ClientId: %s, RedirectURI: %s, AccessType: %s, Scope: %s}", c.Provider, c.ClientId, c.RedirectURI, c.AccessType, c.Scope)
}
