package model

import "fmt"

type LoginResponse struct {
	AccessToken string `json:"auth_token"`
	UserId      string `json:"user_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	LoginType   string `json:"login_type"`
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

type VerifyAdminResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type OAuthProvider struct {
	Provider    string `json:"provider"`
	ClientId    string `json:"client_id"`
	BaseUrl     string `json:"base_url"`
	RedirectURI string `json:"redirect_uri"`
	AccessType  string `json:"access_type"`
	Scope       string `json:"scope"`
}

type OAuthProviderResponse struct {
	Clients    []OAuthProvider `json:"clients"`
	Success    bool            `json:"success"`
	StatusCode int             `json:"status_code"`
}

type OAuthCallbackResponse struct {
	Success   bool   `json:"success"`
	UserId    string `json:"user_id"`
	AuthToken string `json:"auth_token"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	IsNewUser bool   `json:"is_new_user"`
	Message   string `json:"message"`
	LoginType string `json:"login_type"`
}

func (c OAuthProvider) String() string {
	return fmt.Sprintf("OAuthProvider{Provider: %s, ClientId: %s, RedirectURI: %s, AccessType: %s, Scope: %s}", c.Provider, c.ClientId, c.RedirectURI, c.AccessType, c.Scope)
}
