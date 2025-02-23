package oauth_dao

import (
	MySQL "github.com/akgarg0472/urlshortener-auth-service/database"
	"github.com/akgarg0472/urlshortener-auth-service/internal/entity"
	"github.com/akgarg0472/urlshortener-auth-service/internal/logger"
	"go.uber.org/zap"
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
		if logger.IsErrorEnabled() {
			logger.Error("Error fetching oAuth providers", zap.Error(result.Error))
		}
		return oAuthProviders
	}

	return oAuthProviders
}
