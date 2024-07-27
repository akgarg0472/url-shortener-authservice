package oauth_dao

import (
	MySQL "github.com/akgarg0472/urlshortener-auth-service/database"
	"github.com/akgarg0472/urlshortener-auth-service/internal/dao/entity"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
)

var (
	logger = Logger.GetLogger("oAuthDao.go")
)

func FetchOAuthClients() []entity.OAuthClient {
	logger.Info("Fetching OAuth clients from database")

	db := MySQL.GetInstance("", "FetchOAuthClients")

	if db == nil {
		logger.Error("Error getting DB instance")
		panic("Failed to obbtain DB instance")
	}

	var oAuthClients []entity.OAuthClient

	result := db.Find(&oAuthClients)

	if result.Error != nil {
		logger.Error("Error fetching oAuth clients: {}", result.Error)
		panic(result.Error)
	}

	return oAuthClients
}
