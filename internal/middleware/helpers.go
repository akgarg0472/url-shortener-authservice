package middleware

import (
	"encoding/json"
	"net/http"
)

func decodeRequestBody(httpRequest *http.Request, ref interface{}) error {
	return json.NewDecoder(httpRequest.Body).Decode(&ref)
}

func writeErrorResponse(responseWriter http.ResponseWriter, statusCode int, message []byte) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)
	_, err := responseWriter.Write(message)
	if err != nil {
		return
	}
}
