package model

import "fmt"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignupRequest struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name" validate:"required"`
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
	BusinessDetails string `json:"business_details"`
	PhoneNumber     string `json:"phone_number" validate:"required"`
	City            string `json:"city" validate:"required"`
	ZipCode         string `json:"zipcode" validate:"required"`
	Country         string `json:"country" validate:"required"`
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

func (request LoginRequest) String() string {
	return fmt.Sprintf("Email: %s", request.Email)
}

func (request SignupRequest) String() string {
	// TODO: implement method
	return fmt.Sprintf("Email: %s, UserId: %s", request.Email, request.LastName)
}
