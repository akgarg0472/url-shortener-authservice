package model

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/akgarg0472/urlshortener-auth-service/constants"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (r LoginRequest) String() string {
	maskedPassword := maskString(r.Password, true)
	return fmt.Sprintf("{Email: %s, Password: %s}", r.Email, maskedPassword)
}

type SignupRequest struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

func (r SignupRequest) String() string {
	return fmt.Sprintf("{Name: %s, Email: %s, Password: %s, ConfirmPassword: %s}", r.Name, r.Email, maskString(r.Password, true), maskString(r.ConfirmPassword, true))
}

type LogoutRequest struct {
	UserId string `json:"user_id" validate:"required"`
}

func (r LogoutRequest) String() string {
	return fmt.Sprintf("{UserId: %s}", r.UserId)
}

type ValidateTokenRequest struct {
	UserId    string `json:"user_id" validate:"required"`
	AuthToken string `json:"auth_token" validate:"required"`
}

func (r ValidateTokenRequest) String() string {
	return fmt.Sprintf("{UserId: %s}", r.UserId)
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required"`
}

func (r ForgotPasswordRequest) String() string {
	return fmt.Sprintf("{Email: %s}", r.Email)
}

type ResetPasswordRequest struct {
	ResetPasswordToken string `json:"token" validate:"required"`
	Email              string `json:"email" validate:"required"`
	Password           string `json:"password" validate:"required"`
	ConfirmPassword    string `json:"confirm_password" validate:"required"`
}

func (r ResetPasswordRequest) String() string {
	return fmt.Sprintf("{ResetPasswordToken: %s, Email: %s, Password: %s, ConfirmPassword: %s}", maskString(r.ResetPasswordToken, false), r.Email, maskString(r.Password, true), maskString(r.ConfirmPassword, true))
}

type OAuthCallbackRequest struct {
	State    string                  `json:"state"`
	Code     string                  `json:"auth_code"`
	Scope    string                  `json:"scope"`
	Provider constants.OAuthProvider `json:"provider"`
}

func (r OAuthCallbackRequest) String() string {
	return fmt.Sprintf("{State: %s, Code: %s, Scope: %s, Provider: %s}", maskString(r.State, false), maskString(r.Code, true), r.Scope, r.Provider)
}

type VerifyAdminRequest struct {
	UserId string `json:"user_id"`
}

func (r VerifyAdminRequest) String() string {
	return fmt.Sprintf("{UserId: %s}", r.UserId)
}

func maskString(input string, isPassword bool) string {
	if len(input) == 0 {
		return input
	}

	length := len(input)
	maskedArray := []rune(input)

	if isPassword {
		if length <= 2 {
			return input
		}
		for i := 1; i < length-1; i++ {
			maskedArray[i] = '*'
		}
	} else {
		maskCount := length / 2
		maskedIndices := make(map[int]bool)
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))

		for len(maskedIndices) < maskCount {
			index := rng.Intn(length)
			if !maskedIndices[index] {
				maskedArray[index] = '*'
				maskedIndices[index] = true
			}
		}
	}

	return string(maskedArray)
}
