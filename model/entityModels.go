package model

import "time"

type User struct {
	Id                  string
	Email               string
	Password            string
	Scopes              string
	FirstName           string
	LastName            string
	PhoneNumber         string
	City                string
	Country             string
	ZipCode             string
	BusinessDetails     string
	ForgotPasswordToken string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (u *User) String() string {
	return "User [id=" + u.Id + ", email=" + u.Email + ", scopes=" + u.Scopes + ", createdAt=" + u.CreatedAt.String() + ", updatedAt=" + u.UpdatedAt.String() + "]"
}
