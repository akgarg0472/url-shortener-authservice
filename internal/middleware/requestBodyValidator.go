//go:build !staticcheck
// +build !staticcheck

package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	AuthModels "github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	utils "github.com/akgarg0472/urlshortener-auth-service/utils"
)

var logger = Logger.GetLogger("requestBodyValidator.go")

func LoginRequestBodyValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestId := httpRequest.Header.Get("Request-ID")

		var loginRequest AuthModels.LoginRequest

		decodeError := decodeRequestBody(httpRequest, &loginRequest)

		if decodeError != nil {
			logger.Error("[{}]: Error decoding request body: {}", requestId, decodeError.Error())
			resp := generateErrorResponse("Invalid request body", 400)
			writeErrorResponse(responseWriter, http.StatusBadRequest, resp)
			return
		}

		validationErrors := utils.ValidateRequestFields(loginRequest)

		if validationErrors != nil {
			logger.Error("[{}]: Login Request Validation failed: {}", requestId, validationErrors)
			errorResponse, _ := json.Marshal(validationErrors)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		ctx := context.WithValue(httpRequest.Context(), "loginRequest", loginRequest)

		next.ServeHTTP(responseWriter, httpRequest.WithContext(ctx))
	})
}

func SignupRequestBodyValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestId := httpRequest.Header.Get("Request-ID")

		var signupRequest AuthModels.SignupRequest

		decodeError := decodeRequestBody(httpRequest, &signupRequest)

		if decodeError != nil {
			logger.Error("[{}]: Error decoding request body: {}", requestId, decodeError.Error())
			resp := generateErrorResponse("Invalid request body", 400)
			writeErrorResponse(responseWriter, http.StatusBadRequest, resp)
			return
		}

		validationErrors := utils.ValidateRequestFields(signupRequest)

		if validationErrors != nil {
			logger.Error("[{}]: Signup Request Validation failed")
			errorResponse, _ := json.Marshal(validationErrors)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		ctx := context.WithValue(httpRequest.Context(), "signupRequest", signupRequest)

		next.ServeHTTP(responseWriter, httpRequest.WithContext(ctx))
	})
}

func decodeRequestBody(httpRequest *http.Request, ref interface{}) error {
	return json.NewDecoder(httpRequest.Body).Decode(&ref)
}

func writeErrorResponse(responseWriter http.ResponseWriter, statusCode int, message []byte) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)
	responseWriter.Write(message)
}

func generateErrorResponse(message string, errorCode int16) []byte {
	errorResponse := AuthModels.ErrorResponse{
		Message:   message,
		ErrorCode: int(errorCode),
	}
	resp, _ := json.Marshal(errorResponse)
	return resp
}
