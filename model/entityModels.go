package model

import (
	"github.com/akgarg0472/urlshortener-auth-service/constants"
)

type User struct {
	Id                  string
	Name                string
	Email               string
	Password            string
	Scopes              string
	ForgotPasswordToken string
	OAuthId             string
	OAuthProvider       string
	LastLoginAt         int64
	PasswordChangedAt   int64
	IsDeleted           bool
	LoginType           constants.UserEntityLoginType
}

func (u User) String() string {
	return "{id=" + u.Id + ", name=" + u.Name + ", email=" + u.Email + ", oAuthId=" + u.OAuthId + ", OAuthProvider=" + u.OAuthProvider + ", scopes=" + u.Scopes + "}"
}
