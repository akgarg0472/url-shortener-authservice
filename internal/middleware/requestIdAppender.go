package middleware

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func AddRequestIdHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestId := generateRequestID()

		httpRequest.Header.Add("Request-ID", requestId)

		next.ServeHTTP(responseWriter, httpRequest)
	})
}

func generateRequestID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
