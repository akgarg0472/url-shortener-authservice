package middleware

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func AddRequestIdHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := generateRequestID()

		r.Header.Add("Request-ID", requestId)

		next.ServeHTTP(w, r)
	})
}

func generateRequestID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
