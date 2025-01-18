package handler

import (
	"net/http"

	oauthservice "github.com/akgarg0472/urlshortener-auth-service/internal/service/auth/oauth"
	"github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
)

var oauthLogger = Logger.GetLogger("authHandler.go")

func GetOAuthProvidersHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	providers := httpRequest.URL.Query().Get("provider")
	clientIds := oauthservice.GetOAuthProvider(providers)

	response := model.OAuthProviderResponse{
		Clients:    clientIds,
		Success:    true,
		StatusCode: 200,
	}

	sendResponseToClient(responseWriter, "", response, nil, 200)
}

func OAuthCallbackHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get("X-Request-Id")
	oAuthCallbackRequest := context.Value(utils.RequestContextKeys.OAuthCallbackRequestKey).(model.OAuthCallbackRequest)

	oauthLogger.Trace("[{}]: OAuth Callback request received on handler -> {}", requestId, oAuthCallbackRequest)

	oAuthCallbackResponse, oAuthCallbackError := oauthservice.ProcessCallbackRequest(requestId, oAuthCallbackRequest)

	sendResponseToClient(responseWriter, requestId, oAuthCallbackResponse, oAuthCallbackError, 200)
}
