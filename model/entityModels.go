package model

type User struct {
	Id                  string
	Name                string
	Email               string
	Password            string
	Scopes              string
	ForgotPasswordToken string
	LastLoginAt         int64
	PasswordChangedAt   int64
	IsDeleted           bool
}

func (u *User) String() string {
	return "User [id=" + u.Id + ", name=" + u.Name + ", email=" + u.Email + ", scopes=" + u.Scopes + "]"
}
