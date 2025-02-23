package middleware

import (
	"fmt"
	"net/http"

	"github.com/akgarg0472/urlshortener-auth-service/constants"
	"github.com/akgarg0472/urlshortener-auth-service/internal/logger"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
	"go.uber.org/zap"
)

func ValidateRequestJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestId := httpRequest.Header.Get(constants.RequestIdHeaderName)
		contentTypeHeader := "Content-Type"
		applicationJsonContentTypeHeader := "application/json"

		contentType := httpRequest.Header.Get(contentTypeHeader)

		if contentType == "" {
			if logger.IsErrorEnabled() {
				logger.Error("Content-Type is missing",
					zap.String(constants.RequestIdLogKey, requestId),
				)
			}
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
			if logger.IsErrorEnabled() {
				logger.Error("Content-Type not supported",
					zap.String(constants.RequestIdLogKey, requestId),
					zap.String("Content-Type", contentType),
				)
			}
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
