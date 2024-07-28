package middleware

import (
	"fmt"
	"net/http"

	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	utils "github.com/akgarg0472/urlshortener-auth-service/utils"
)

var (
	httpRequestHeaderValidatorLogger = Logger.GetLogger("httpRequestHeaderValidator.go")
)

func ValidateRequestJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestId := httpRequest.Header.Get("Request-ID")
		contentTypeHeader := "Content-Type"
		applicationJsonContentTypeHeader := "application/json"

		contentType := httpRequest.Header.Get(contentTypeHeader)

		if contentType == "" {
			httpRequestHeaderValidatorLogger.Error("[{}]: Content-Type is missing", requestId)
			errorResponseJson := utils.GetErrorResponseByte("Content-Type is missing", 400)
			responseWriter.Header().Set(contentTypeHeader, applicationJsonContentTypeHeader)
			responseWriter.WriteHeader(http.StatusBadRequest)
			_, err := responseWriter.Write(errorResponseJson)
			if err != nil {
				return
			}
			return
		}

		if contentType != applicationJsonContentTypeHeader {
			httpRequestHeaderValidatorLogger.Error("[{}]: Content-Type '%s' not supported", requestId, contentType)
			errorResponseJson := utils.GetErrorResponseByte(fmt.Sprintf("Content-Type '%s' not supported", contentType), 400)
			responseWriter.Header().Set(contentTypeHeader, applicationJsonContentTypeHeader)
			responseWriter.WriteHeader(http.StatusBadRequest)
			_, err := responseWriter.Write(errorResponseJson)
			if err != nil {
				return
			}
			return
		}

		next.ServeHTTP(responseWriter, httpRequest)
	})
}
