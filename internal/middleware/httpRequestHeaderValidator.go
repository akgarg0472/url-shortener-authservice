package middleware

import (
	"fmt"
	"net/http"

	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	utils "github.com/akgarg0472/urlshortener-auth-service/utils"
)

var (
	hrhvLogger = Logger.GetLogger("httpRequestHeaderValidator.go")
)

func ValidateRequestJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestId := httpRequest.Header.Get("Request-ID")
		contentType := httpRequest.Header.Get("Content-Type")

		if contentType == "" {
			hrhvLogger.Error("[{}]: Content-Type is missing", requestId)
			errorResponseJson := utils.GetErrorResponseByte("Content-Type is missing", 400)
			responseWriter.Header().Set("Content-Type", "application/json")
			responseWriter.WriteHeader(http.StatusBadRequest)
			responseWriter.Write(errorResponseJson)
			return
		}

		if contentType != "application/json" {
			hrhvLogger.Error("[{}]: Content-Type '%s' not supported", requestId, contentType)
			errorResponseJson := utils.GetErrorResponseByte(fmt.Sprintf("Content-Type '%s' not supported", contentType), 400)
			responseWriter.Header().Set("Content-Type", "application/json")
			responseWriter.WriteHeader(http.StatusBadRequest)
			responseWriter.Write(errorResponseJson)
			return
		}

		next.ServeHTTP(responseWriter, httpRequest)
	})
}
