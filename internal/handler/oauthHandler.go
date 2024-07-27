package handler

import (
	"net/http"
	"strings"

	authModels "github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
)

var oauthLogger = Logger.GetLogger("oauthHandler.go")

// Handler Function to handle login request
func GetOAuthClients(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	requestId := httpRequest.Header.Get("Request-ID")

	provider := httpRequest.URL.Query().Get("provider")
	providers := strings.Split(provider, ",")
	oauthLogger.Info("Fetching oauth client for OAuth provider(s): {}", providers)

	clientIds := []authModels.OAuthClient{
		{
			Provider:    "google",
			ClientId:    "186092380904-4e3sen6d26df4tvqaj8ou4vbvkfbjt47.apps.googleusercontent.com",
			RedirectURI: "http://localhost:3000/oauth/google/success",
			AccessType:  "offline",
			Scope:       "profile email",
		},
	}

	oauthLogger.Info("Sending OAuth clients: {}", clientIds)

	sendResponseToClient(responseWriter, requestId, clientIds, nil, 200)
}
