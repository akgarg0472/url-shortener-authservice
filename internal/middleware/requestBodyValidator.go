package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	AuthModels "github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	utils "github.com/akgarg0472/urlshortener-auth-service/utils"
)

var rbvLogger = Logger.GetLogger("requestBodyValidator.go")

func LoginRequestBodyValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestId := httpRequest.Header.Get("Request-ID")

		var loginRequest AuthModels.LoginRequest

		decodeError := decodeRequestBody(httpRequest, &loginRequest)

		if decodeError != nil {
			rbvLogger.Error("[{}]: Error decoding login request body: {}", requestId, decodeError.Error())
			resp := utils.GetErrorResponse("Invalid request body", 400)
			errorResponseJson, _ := utils.ConvertToJsonBytes(resp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponseJson)
			return
		}

		validationErrors := utils.ValidateRequestFields(loginRequest)

		if validationErrors != nil {
			rbvLogger.Error("[{}]: Login Request Validation failed: {}", requestId, validationErrors)
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
			rbvLogger.Error("[{}]: Error decoding signup request body: {}", requestId, decodeError.Error())
			resp := utils.GetErrorResponse("Invalid request body", 400)
			errorJsonResponse, _ := utils.ConvertToJsonBytes(resp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorJsonResponse)
			return
		}

		validationErrors := utils.ValidateRequestFields(signupRequest)

		if validationErrors != nil {
			rbvLogger.Error("[{}]: Signup Request Validation failed")
			errorResponse, _ := json.Marshal(validationErrors)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		ctx := context.WithValue(httpRequest.Context(), "signupRequest", signupRequest)

		next.ServeHTTP(responseWriter, httpRequest.WithContext(ctx))
	})
}

func LogoutRequestBodyValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestId := httpRequest.Header.Get("Request-ID")

		var logoutRequest AuthModels.LogoutRequest

		decodeError := decodeRequestBody(httpRequest, &logoutRequest)

		if decodeError != nil {
			rbvLogger.Error("[{}]: Error decoding logout request body: {}", requestId, decodeError.Error())
			resp := utils.GetErrorResponse("Invalid request body", 400)
			errorJsonResponse, _ := utils.ConvertToJsonBytes(resp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorJsonResponse)
			return
		}

		validationErrors := utils.ValidateRequestFields(logoutRequest)

		if validationErrors != nil {
			rbvLogger.Error("[{}]: Logout Request Validation failed")
			errorResponse, _ := json.Marshal(validationErrors)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		ctx := context.WithValue(httpRequest.Context(), "logoutRequest", logoutRequest)

		next.ServeHTTP(responseWriter, httpRequest.WithContext(ctx))
	})
}

func VerifyTokenRequestBodyValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestId := httpRequest.Header.Get("Request-ID")

		var validateTokenRequest AuthModels.ValidateTokenRequest

		decodeError := decodeRequestBody(httpRequest, &validateTokenRequest)

		if decodeError != nil {
			rbvLogger.Error("[{}]: Error decoding validate token request body: {}", requestId, decodeError.Error())
			resp := utils.GetErrorResponse("Invalid request body", 400)
			errorJsonResponse, _ := utils.ConvertToJsonBytes(resp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorJsonResponse)
			return
		}

		validationErrors := utils.ValidateRequestFields(validateTokenRequest)

		if validationErrors != nil {
			rbvLogger.Error("[{}]: Validate Token Request Validation failed")
			errorResponse, _ := json.Marshal(validationErrors)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		ctx := context.WithValue(httpRequest.Context(), "validateTokenRequest", validateTokenRequest)

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
