package model

import "fmt"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignupRequest struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

type LogoutRequest struct {
	AuthToken string `json:"auth_token" validate:"required"`
	UserId    string `json:"user_id" validate:"required"`
}

type ValidateTokenRequest struct {
	AuthToken string `json:"auth_token" validate:"required"`
	UserId    string `json:"user_id" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required"`
}

type ResetPasswordRequest struct {
	ResetPasswordToken string `json:"token" validate:"required"`
	Email              string `json:"email" validate:"required"`
	Password           string `json:"password" validate:"required"`
	ConfirmPassword    string `json:"confirm_password" validate:"required"`
}

type OAuthProvider string

const (
	OAUTH_PROVIDER_GOOGLE string = "google"
	OAUTH_PROVIDER_GITHUB string = "github"
)

type OAuthCallbackRequest struct {
	State    string        `json:"state"`
	Code     string        `json:"auth_code"`
	Scope    string        `json:"scope"`
	Provider OAuthProvider `json:"provider"`
}

func (request LoginRequest) String() string {
	return fmt.Sprintf("Email: %s", request.Email)
}

func (request SignupRequest) String() string {
	return fmt.Sprintf("Email: %s, UserId: %s", request.Email, request.Name)
}

func (r OAuthCallbackRequest) String() string {
	return fmt.Sprintf("OAuthCallbackRequest {State: %s, Code: %s, Scope: %s, Provider: %s}", r.State, r.Code, r.Scope, r.Provider)
}
