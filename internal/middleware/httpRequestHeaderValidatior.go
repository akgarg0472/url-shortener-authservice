package middleware

import (
	"net/http"
)

func ValidateRequestContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		contentType := httpRequest.Header.Get("Content-Type")

		if contentType == "" {
			responseWriter.WriteHeader(http.StatusBadRequest)
			responseWriter.Write([]byte("Content-Type is invalid"))
			return
		}

		if contentType != "application/json" {
			responseWriter.WriteHeader(http.StatusBadRequest)
			responseWriter.Write([]byte("Content-Type '" + contentType + "' not supported"))
			return
		}

		next.ServeHTTP(responseWriter, httpRequest)
	})
}
