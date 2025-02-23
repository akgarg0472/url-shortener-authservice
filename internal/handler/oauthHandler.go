package handler

import (
	"net/http"

	"github.com/akgarg0472/urlshortener-auth-service/constants"
	"github.com/akgarg0472/urlshortener-auth-service/internal/logger"
	oauth_service "github.com/akgarg0472/urlshortener-auth-service/internal/service/auth/oauth"
	"github.com/akgarg0472/urlshortener-auth-service/model"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
	"go.uber.org/zap"
)

func GetOAuthProvidersHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	providers := httpRequest.URL.Query().Get("provider")
	clientIds := oauth_service.GetOAuthProvider(providers)

	response := model.OAuthProviderResponse{
		Clients:    clientIds,
		Success:    true,
		StatusCode: 200,
	}

	sendResponseToClient(responseWriter, "", response, nil, 200)
}

func OAuthCallbackHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get(constants.RequestIdHeaderName)
	oAuthCallbackRequest := context.Value(utils.RequestContextKeys.OAuthCallbackRequestKey).(model.OAuthCallbackRequest)

	if logger.IsDebugEnabled() {
		logger.Debug("OAuth Callback request received on handler",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.Any("oAuthCallbackRequest", oAuthCallbackRequest),
		)
	}

	oAuthCallbackResponse, oAuthCallbackError := oauth_service.ProcessCallbackRequest(requestId, oAuthCallbackRequest)

	sendResponseToClient(responseWriter, requestId, oAuthCallbackResponse, oAuthCallbackError, 200)
}
