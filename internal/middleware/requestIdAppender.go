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
	u := uuid.New().String()
	u = strings.ReplaceAll(u, "-", "")
	if len(u) > 16 {
		u = u[:16]
	}
	return u
}
