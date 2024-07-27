package entity

type OAuthClient struct {
	ID          uint8  `gorm:"primaryKey"`
	Provider    string `gorm:"size:255;not null;unique"`
	ClientID    string `gorm:"size:255;unique;not null"`
	RedirectURI string `gorm:"type:text;not null"`
	AccessType  string `gorm:"size:50"`
	Scope       string `gorm:"type:text"`
}

func (OAuthClient) TableName() string {
	return "oauth_clients"
}
