package middleware

import (
	"net/http"
	"strings"

	"github.com/akgarg0472/urlshortener-auth-service/constants"
	"github.com/google/uuid"
)

func AddRequestIdHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		if httpRequest.Header.Get(constants.RequestIdHeaderName) == "" {
			httpRequest = httpRequest.Clone(httpRequest.Context())
			httpRequest.Header.Set(constants.RequestIdHeaderName, generateRequestID())
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
