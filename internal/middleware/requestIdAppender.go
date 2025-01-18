package middleware

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
)

var requestIdHeader = "X-Request-Id"

func AddRequestIdHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		if httpRequest.Header.Get(requestIdHeader) == "" {
			httpRequest = httpRequest.Clone(httpRequest.Context())
			httpRequest.Header.Set(requestIdHeader, generateRequestID())
		}

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
