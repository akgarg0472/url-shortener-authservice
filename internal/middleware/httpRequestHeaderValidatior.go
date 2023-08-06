package middleware

import (
	"net/http"
)

func ValidateRequestContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if contentType == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Content-Type is invalid"))
			return
		}

		if contentType != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Content-Type '" + contentType + "' not supported"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
