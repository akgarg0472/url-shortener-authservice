package model

import "fmt"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignupRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (request LoginRequest) String() string {
	return fmt.Sprintf("Email: %s", request.Email)
}

func (request SignupRequest) String() string {
	return fmt.Sprintf("Email: %s", request.Email)
}
