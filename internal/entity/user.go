package entity

import "github.com/akgarg0472/urlshortener-auth-service/constants"

type User struct {
	Id                    string                    `gorm:"primaryKey;size:128" json:"id"`                 // varchar(128)
	Email                 *string                   `gorm:"uniqueIndex;size:255" json:"email"`             // varchar(255)
	Password              *string                   `gorm:"size:255;" json:"password"`                     // varchar(255)
	Scopes                string                    `gorm:"size:32" json:"scopes"`                         // varchar(32)
	Name                  string                    `gorm:"size:255" json:"name"`                          // varchar(255)
	Bio                   *string                   `gorm:"type:text" json:"bio,omitempty"`                // text
	ProfilePictureURL     *string                   `gorm:"size:255" json:"profile_picture_url,omitempty"` // varchar(255)
	Phone                 *string                   `gorm:"size:20" json:"phone,omitempty"`                // varchar(20)
	UserLoginType         enums.UserEntityLoginType `gorm:"type:varchar(50)" json:"login_type"`
	OAuthId               *string                   `gorm:"column:oauth_id;uniqueIndex;size:255" json:"oauth_id,omitempty"`
	OAuthProvider         *string                   `gorm:"column:oauth_provider;size:16" json:"oauth_provider,omitempty"`
	PremiumAccount        bool                      `gorm:"default:0" json:"premium_account"`                      // tinyint(1)
	City                  *string                   `gorm:"size:50" json:"city,omitempty"`                         // varchar(50)
	State                 *string                   `gorm:"size:50" json:"state,omitempty"`                        // varchar(50)
	Country               *string                   `gorm:"size:50" json:"country,omitempty"`                      // varchar(50)
	Zipcode               *string                   `gorm:"size:16" json:"zipcode,omitempty"`                      // varchar(16)
	BusinessDetails       *string                   `gorm:"type:text" json:"business_details,omitempty"`           // text
	ForgotPasswordToken   *string                   `gorm:"size:255" json:"forgot_password_token,omitempty"`       // varchar(255)
	LastPasswordChangedAt *int64                    `gorm:"type:bigint" json:"last_password_changed_at,omitempty"` // bigint
	LastLoginAt           *int64                    `gorm:"type:bigint" json:"last_login_at,omitempty"`            // bigint
	IsDeleted             bool                      `gorm:"default:0" json:"is_deleted"`                           // tinyint(1)
	CreatedAt             int64                     `gorm:"type:bigint;" json:"created_at"`                        // timestamp
	UpdatedAt             int64                     `gorm:"type:bigint;" json:"updated_at"`                        // timestamp
}

func (User) TableName() string {
	return "users"
}
