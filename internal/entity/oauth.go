package entity

type OAuthProvider struct {
	ID          uint8  `gorm:"primaryKey"`
	Provider    string `gorm:"size:255;not null;unique"`
	ClientID    string `gorm:"size:255;unique;not null"`
	BaseUrl     string `gorm:"size:255;unique;not null"`
	RedirectURI string `gorm:"type:text;not null"`
	AccessType  string `gorm:"size:50"`
	Scope       string `gorm:"type:text"`
}

func (OAuthProvider) TableName() string {
	return "oauth_providers"
}
