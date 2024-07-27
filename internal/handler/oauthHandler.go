package handler

import (
	"net/http"

	oAuthService "github.com/akgarg0472/urlshortener-auth-service/internal/service/auth/oauth"
	"github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
)

var oauthLogger = Logger.GetLogger("authHandler.go")

// Handler Function to handle login request
func GetOAuthClientHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	providers := httpRequest.URL.Query().Get("provider")
	clientIds := oAuthService.GetOAuthClient(providers)

	response := model.OAuthClientResponse{
		Clients:    clientIds,
		Success:    true,
		StatusCode: 200,
	}

	sendResponseToClient(responseWriter, "", response, nil, 200)
}

func OAuthCallbackHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get("Request-ID")
	oAuthCallbackRequest := context.Value("oAuthCallbackRequest").(model.OAuthCallbackRequest)

	oauthLogger.Trace("[{}]: OAuth Callback request received on handler -> {}", requestId, oAuthCallbackRequest)

	oAuthCallbackResponse, oAuthCallbackError := oAuthService.ProcessCallbackRequest(requestId, oAuthCallbackRequest)

	sendResponseToClient(responseWriter, requestId, oAuthCallbackResponse, oAuthCallbackError, 200)
}
