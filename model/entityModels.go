package model

import "time"

type User struct {
	Id                  string
	Name                string
	Email               string
	Password            string
	Scopes              string
	ForgotPasswordToken string
	LastLoginAt         time.Time
	PasswordChangedAt   time.Time
}

func (u *User) String() string {
	return "User [id=" + u.Id + ", name=" + u.Name + ", email=" + u.Email + ", scopes=" + u.Scopes + "]"
}
