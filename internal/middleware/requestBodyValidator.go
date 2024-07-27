package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

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
			errResp := AuthModels.ErrorResponse{
				Message:   "Request validation failed",
				ErrorCode: 400,
				Errors:    validationErrors,
			}
			errorResponse, _ := json.Marshal(errResp)
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
			errResp := AuthModels.ErrorResponse{
				Message:   "Request validation failed",
				ErrorCode: 400,
				Errors:    validationErrors,
			}
			errorResponse, _ := json.Marshal(errResp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		if strings.TrimSpace(signupRequest.Password) != strings.TrimSpace(signupRequest.ConfirmPassword) {
			rbvLogger.Error("[{}]: Signup Request Validation failed. Passwords mismatch")
			errResp := AuthModels.ErrorResponse{
				Message:   "Request validation failed",
				ErrorCode: 400,
				Errors:    "Password and confirm password mismatch",
			}
			errorResponse, _ := json.Marshal(errResp)
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
			errResp := AuthModels.ErrorResponse{
				Message:   "Request validation failed",
				ErrorCode: 400,
				Errors:    validationErrors,
			}
			errorResponse, _ := json.Marshal(errResp)
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
			errResp := AuthModels.ErrorResponse{
				Message:   "Request validation failed",
				ErrorCode: 400,
				Errors:    validationErrors,
			}
			errorResponse, _ := json.Marshal(errResp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		ctx := context.WithValue(httpRequest.Context(), "validateTokenRequest", validateTokenRequest)

		next.ServeHTTP(responseWriter, httpRequest.WithContext(ctx))
	})
}

func ForgotPasswordRequestBodyValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestId := httpRequest.Header.Get("Request-ID")

		var forgotPasswordRequest AuthModels.ForgotPasswordRequest

		decodeError := decodeRequestBody(httpRequest, &forgotPasswordRequest)

		if decodeError != nil {
			rbvLogger.Error("[{}]: Error decoding forgot password request body: {}", requestId, decodeError.Error())
			resp := utils.GetErrorResponse("Invalid request body", 400)
			errorJsonResponse, _ := utils.ConvertToJsonBytes(resp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorJsonResponse)
			return
		}

		validationErrors := utils.ValidateRequestFields(forgotPasswordRequest)

		if validationErrors != nil {
			rbvLogger.Error("[{}]: Forgot Password Request Validation failed")
			errResp := AuthModels.ErrorResponse{
				Message:   "Request validation failed",
				ErrorCode: 400,
				Errors:    validationErrors,
			}
			errorResponse, _ := json.Marshal(errResp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		ctx := context.WithValue(httpRequest.Context(), "forgotPasswordRequest", forgotPasswordRequest)

		next.ServeHTTP(responseWriter, httpRequest.WithContext(ctx))
	})
}

func ResetPasswordRequestBodyValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestId := httpRequest.Header.Get("Request-ID")

		var resetPasswordRequest AuthModels.ResetPasswordRequest

		decodeError := decodeRequestBody(httpRequest, &resetPasswordRequest)

		if decodeError != nil {
			rbvLogger.Error("[{}]: Error decoding reset password request body: {}", requestId, decodeError.Error())
			resp := utils.GetErrorResponse("Invalid request body", 400)
			errorJsonResponse, _ := utils.ConvertToJsonBytes(resp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorJsonResponse)
			return
		}

		validationErrors := utils.ValidateRequestFields(resetPasswordRequest)

		if validationErrors != nil {
			rbvLogger.Error("[{}]: Reset Password Request Validation failed")
			errResp := AuthModels.ErrorResponse{
				Message:   "Request validation failed",
				ErrorCode: 400,
				Errors:    validationErrors,
			}
			errorResponse, _ := json.Marshal(errResp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		if strings.TrimSpace(resetPasswordRequest.Password) != strings.TrimSpace(resetPasswordRequest.ConfirmPassword) {
			rbvLogger.Error("[{}]: Reset Password Request Validation failed. Passwords mismatch")
			errResp := AuthModels.ErrorResponse{
				Message:   "Request validation failed",
				ErrorCode: 400,
				Errors:    "Password and confirm password mismatch",
			}
			errorResponse, _ := json.Marshal(errResp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		ctx := context.WithValue(httpRequest.Context(), "resetPasswordRequest", resetPasswordRequest)

		next.ServeHTTP(responseWriter, httpRequest.WithContext(ctx))
	})
}

func OAuthCallbackRequestBodyValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestId := httpRequest.Header.Get("Request-ID")

		var oAuthCallbackRequest AuthModels.OAuthCallbackRequest

		decodeError := decodeRequestBody(httpRequest, &oAuthCallbackRequest)

		if decodeError != nil {
			rbvLogger.Error("[{}]: Error decoding oAuth callback request body: {}", requestId, decodeError.Error())
			resp := utils.GetErrorResponse("Invalid request body", 400)
			errorJsonResponse, _ := utils.ConvertToJsonBytes(resp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorJsonResponse)
			return
		}

		validationErrors := utils.ValidateRequestFields(oAuthCallbackRequest)

		if validationErrors != nil {
			rbvLogger.Error("[{}]: OAuth Callback Request Validation failed", requestId)
			errResp := AuthModels.ErrorResponse{
				Message:   "Request validation failed",
				ErrorCode: 400,
				Errors:    validationErrors,
			}
			errorResponse, _ := json.Marshal(errResp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		if oAuthCallbackRequest.Code == "" {
			rbvLogger.Error("[{}]: OAuth Callback Request Validation failed", requestId)
			errResp := AuthModels.ErrorResponse{
				Message:   "Please provide valid auth_code",
				ErrorCode: 400,
				Errors:    validationErrors,
			}
			errorResponse, _ := json.Marshal(errResp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		ctx := context.WithValue(httpRequest.Context(), "oAuthCallbackRequest", oAuthCallbackRequest)

		next.ServeHTTP(responseWriter, httpRequest.WithContext(ctx))
	})
}
