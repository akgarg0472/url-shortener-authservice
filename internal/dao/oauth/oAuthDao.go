package oauth_dao

import (
	MySQL "github.com/akgarg0472/urlshortener-auth-service/database"
	"github.com/akgarg0472/urlshortener-auth-service/internal/entity"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
)

var (
	logger = Logger.GetLogger("oAuthDao.go")
)

func FetchOAuthProviders() []entity.OAuthProvider {
	logger.Info("Fetching OAuth providers from database")

	db := MySQL.GetInstance("", "FetchOAuthProviders")

	if db == nil {
		logger.Error("Error getting DB instance")
		panic("Failed to obtain DB instance")
	}

	var oAuthProviders []entity.OAuthProvider

	result := db.Find(&oAuthProviders)

	if result.Error != nil {
		logger.Error("Error fetching oAuth providers: {}", result.Error)
		panic(result.Error)
	}

	return oAuthProviders
}
