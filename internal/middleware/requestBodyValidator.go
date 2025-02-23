package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/akgarg0472/urlshortener-auth-service/constants"
	"github.com/akgarg0472/urlshortener-auth-service/internal/logger"
	AuthModels "github.com/akgarg0472/urlshortener-auth-service/model"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
	"go.uber.org/zap"
)

var invalidRequestBodyMessage = "Invalid request body"
var requestValidationFailedMessage = "Request validation failed"

func LoginRequestBodyValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestId := httpRequest.Header.Get(constants.RequestIdHeaderName)

		var loginRequest AuthModels.LoginRequest

		decodeError := decodeRequestBody(httpRequest, &loginRequest)

		if decodeError != nil {
			if logger.IsErrorEnabled() {
				logger.Error("Error decoding login request body",
					zap.String(constants.RequestIdLogKey, requestId),
					zap.Error(decodeError),
				)
			}
			resp := utils.GetErrorResponse(invalidRequestBodyMessage, 400)
			errorResponseJson, _ := utils.ConvertToJsonBytes(resp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponseJson)
			return
		}

		validationErrors := utils.ValidateRequestFields(loginRequest)

		if validationErrors != nil {
			if logger.IsErrorEnabled() {
				logger.Error("Login Request Validation failed",
					zap.String(constants.RequestIdLogKey, requestId),
					zap.Any("validation_errors", validationErrors),
				)
			}
			errResp := AuthModels.ErrorResponse{
				Message:   requestValidationFailedMessage,
				ErrorCode: 400,
				Errors:    validationErrors,
			}
			errorResponse, _ := json.Marshal(errResp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		ctx := context.WithValue(httpRequest.Context(), utils.RequestContextKeys.LoginRequestKey, loginRequest)

		next.ServeHTTP(responseWriter, httpRequest.WithContext(ctx))
	})
}

func SignupRequestBodyValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestId := httpRequest.Header.Get(constants.RequestIdHeaderName)

		var signupRequest AuthModels.SignupRequest

		decodeError := decodeRequestBody(httpRequest, &signupRequest)

		if decodeError != nil {
			if logger.IsErrorEnabled() {
				logger.Error("Error decoding signup request body",
					zap.String(constants.RequestIdLogKey, requestId),
					zap.String("error", decodeError.Error()),
				)
			}
			resp := utils.GetErrorResponse(invalidRequestBodyMessage, 400)
			errorJsonResponse, _ := utils.ConvertToJsonBytes(resp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorJsonResponse)
			return
		}

		validationErrors := utils.ValidateRequestFields(signupRequest)

		if validationErrors != nil {
			if logger.IsErrorEnabled() {
				logger.Error("Signup Request Validation failed",
					zap.String(constants.RequestIdLogKey, requestId),
				)
			}
			errResp := AuthModels.ErrorResponse{
				Message:   requestValidationFailedMessage,
				ErrorCode: 400,
				Errors:    validationErrors,
			}
			errorResponse, _ := json.Marshal(errResp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		if strings.TrimSpace(signupRequest.Password) != strings.TrimSpace(signupRequest.ConfirmPassword) {
			if logger.IsErrorEnabled() {
				logger.Error("Signup Request Validation failed. Passwords mismatch",
					zap.String(constants.RequestIdLogKey, requestId),
				)
			}
			errResp := AuthModels.ErrorResponse{
				Message:   requestValidationFailedMessage,
				ErrorCode: 400,
				Errors:    "Password and confirm password mismatch",
			}
			errorResponse, _ := json.Marshal(errResp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		ctx := context.WithValue(httpRequest.Context(), utils.RequestContextKeys.SignupRequestKey, signupRequest)

		next.ServeHTTP(responseWriter, httpRequest.WithContext(ctx))
	})
}

func LogoutRequestBodyValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestId := httpRequest.Header.Get(constants.RequestIdHeaderName)

		var logoutRequest AuthModels.LogoutRequest

		decodeError := decodeRequestBody(httpRequest, &logoutRequest)

		if decodeError != nil {
			if logger.IsErrorEnabled() {
				logger.Error("Error decoding logout request body",
					zap.String(constants.RequestIdLogKey, requestId),
					zap.String("error", decodeError.Error()),
				)
			}
			resp := utils.GetErrorResponse(invalidRequestBodyMessage, 400)
			errorJsonResponse, _ := utils.ConvertToJsonBytes(resp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorJsonResponse)
			return
		}

		validationErrors := utils.ValidateRequestFields(logoutRequest)

		if validationErrors != nil {
			logger.Error(
				"Logout Request Validation failed",
				zap.String(constants.RequestIdLogKey, requestId),
			)
			errResp := AuthModels.ErrorResponse{
				Message:   requestValidationFailedMessage,
				ErrorCode: 400,
				Errors:    validationErrors,
			}
			errorResponse, _ := json.Marshal(errResp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		ctx := context.WithValue(httpRequest.Context(), utils.RequestContextKeys.LogoutRequestKey, logoutRequest)

		next.ServeHTTP(responseWriter, httpRequest.WithContext(ctx))
	})
}

func VerifyTokenRequestBodyValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestId := httpRequest.Header.Get(constants.RequestIdHeaderName)

		var validateTokenRequest AuthModels.ValidateTokenRequest

		decodeError := decodeRequestBody(httpRequest, &validateTokenRequest)

		if decodeError != nil {
			if logger.IsErrorEnabled() {
				logger.Error("Error decoding validate token request body",
					zap.String(constants.RequestIdLogKey, requestId),
					zap.String("error", decodeError.Error()),
				)
			}
			resp := utils.GetErrorResponse(invalidRequestBodyMessage, 400)
			errorJsonResponse, _ := utils.ConvertToJsonBytes(resp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorJsonResponse)
			return
		}

		validationErrors := utils.ValidateRequestFields(validateTokenRequest)

		if validationErrors != nil {
			if logger.IsErrorEnabled() {
				logger.Error("Validate Token Request Validation failed",
					zap.String(constants.RequestIdLogKey, requestId),
				)
			}
			errResp := AuthModels.ErrorResponse{
				Message:   requestValidationFailedMessage,
				ErrorCode: 400,
				Errors:    validationErrors,
			}
			errorResponse, _ := json.Marshal(errResp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		ctx := context.WithValue(httpRequest.Context(), utils.RequestContextKeys.ValidateTokenRequestKey, validateTokenRequest)

		next.ServeHTTP(responseWriter, httpRequest.WithContext(ctx))
	})
}

func ForgotPasswordRequestBodyValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestId := httpRequest.Header.Get(constants.RequestIdHeaderName)

		var forgotPasswordRequest AuthModels.ForgotPasswordRequest

		decodeError := decodeRequestBody(httpRequest, &forgotPasswordRequest)

		if decodeError != nil {
			if logger.IsErrorEnabled() {
				logger.Error("Error decoding forgot password request body",
					zap.String(constants.RequestIdLogKey, requestId),
					zap.Error(decodeError),
				)
			}
			resp := utils.GetErrorResponse(invalidRequestBodyMessage, 400)
			errorJsonResponse, _ := utils.ConvertToJsonBytes(resp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorJsonResponse)
			return
		}

		validationErrors := utils.ValidateRequestFields(forgotPasswordRequest)

		if validationErrors != nil {
			if logger.IsErrorEnabled() {
				logger.Error("Forgot Password Request Validation failed",
					zap.String(constants.RequestIdLogKey, requestId),
				)
			}
			errResp := AuthModels.ErrorResponse{
				Message:   requestValidationFailedMessage,
				ErrorCode: 400,
				Errors:    validationErrors,
			}
			errorResponse, _ := json.Marshal(errResp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		ctx := context.WithValue(httpRequest.Context(), utils.RequestContextKeys.ForgotPasswordRequestKey, forgotPasswordRequest)

		next.ServeHTTP(responseWriter, httpRequest.WithContext(ctx))
	})
}

func ResetPasswordRequestBodyValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestId := httpRequest.Header.Get(constants.RequestIdHeaderName)

		var resetPasswordRequest AuthModels.ResetPasswordRequest

		decodeError := decodeRequestBody(httpRequest, &resetPasswordRequest)

		if decodeError != nil {
			if logger.IsErrorEnabled() {
				logger.Error("Error decoding reset password request body",
					zap.String(constants.RequestIdLogKey, requestId),
					zap.Error(decodeError),
				)
			}
			resp := utils.GetErrorResponse(invalidRequestBodyMessage, 400)
			errorJsonResponse, _ := utils.ConvertToJsonBytes(resp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorJsonResponse)
			return
		}

		validationErrors := utils.ValidateRequestFields(resetPasswordRequest)

		if validationErrors != nil {
			if logger.IsErrorEnabled() {
				logger.Error("Reset Password Request Validation failed",
					zap.String(constants.RequestIdLogKey, requestId),
				)
			}
			errResp := AuthModels.ErrorResponse{
				Message:   requestValidationFailedMessage,
				ErrorCode: 400,
				Errors:    validationErrors,
			}
			errorResponse, _ := json.Marshal(errResp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		if strings.TrimSpace(resetPasswordRequest.Password) != strings.TrimSpace(resetPasswordRequest.ConfirmPassword) {
			if logger.IsErrorEnabled() {
				logger.Error("Reset Password Request Validation failed. Passwords mismatch",
					zap.String(constants.RequestIdLogKey, requestId),
				)
			}
			errResp := AuthModels.ErrorResponse{
				Message:   requestValidationFailedMessage,
				ErrorCode: 400,
				Errors:    "Password and confirm password mismatch",
			}
			errorResponse, _ := json.Marshal(errResp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		ctx := context.WithValue(httpRequest.Context(), utils.RequestContextKeys.ResetPasswordRequestKey, resetPasswordRequest)

		next.ServeHTTP(responseWriter, httpRequest.WithContext(ctx))
	})
}

func OAuthCallbackRequestBodyValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestId := httpRequest.Header.Get(constants.RequestIdHeaderName)

		var oAuthCallbackRequest AuthModels.OAuthCallbackRequest

		decodeError := decodeRequestBody(httpRequest, &oAuthCallbackRequest)

		if decodeError != nil {
			if logger.IsErrorEnabled() {
				logger.Error("Error decoding oAuth callback request body",
					zap.String(constants.RequestIdLogKey, requestId),
					zap.Error(decodeError),
				)
			}
			resp := utils.GetErrorResponse(invalidRequestBodyMessage, 400)
			errorJsonResponse, _ := utils.ConvertToJsonBytes(resp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorJsonResponse)
			return
		}

		validationErrors := utils.ValidateRequestFields(oAuthCallbackRequest)

		if validationErrors != nil {
			if logger.IsErrorEnabled() {
				logger.Error("OAuth Callback Request Validation failed",
					zap.String(constants.RequestIdLogKey, requestId),
				)
			}
			errResp := AuthModels.ErrorResponse{
				Message:   requestValidationFailedMessage,
				ErrorCode: 400,
				Errors:    validationErrors,
			}
			errorResponse, _ := json.Marshal(errResp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		if oAuthCallbackRequest.Code == "" {
			if logger.IsErrorEnabled() {
				logger.Error("OAuth Callback Request Validation failed",
					zap.String(constants.RequestIdLogKey, requestId),
				)
			}
			errResp := AuthModels.ErrorResponse{
				Message:   "Please provide valid auth_code",
				ErrorCode: 400,
				Errors:    validationErrors,
			}
			errorResponse, _ := json.Marshal(errResp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		ctx := context.WithValue(httpRequest.Context(), utils.RequestContextKeys.OAuthCallbackRequestKey, oAuthCallbackRequest)

		next.ServeHTTP(responseWriter, httpRequest.WithContext(ctx))
	})
}

func VerifyAdminRequestBodyHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestId := httpRequest.Header.Get(constants.RequestIdHeaderName)

		var verifyAdminRequest AuthModels.VerifyAdminRequest

		decodeError := decodeRequestBody(httpRequest, &verifyAdminRequest)

		if decodeError != nil {
			if logger.IsErrorEnabled() {
				logger.Error("Error decoding verify admin request body",
					zap.String(constants.RequestIdLogKey, requestId),
					zap.Error(decodeError),
				)
			}
			resp := utils.GetErrorResponse(invalidRequestBodyMessage, 400)
			errorJsonResponse, _ := utils.ConvertToJsonBytes(resp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorJsonResponse)
			return
		}

		validationErrors := utils.ValidateRequestFields(verifyAdminRequest)

		if validationErrors != nil {
			if logger.IsErrorEnabled() {
				logger.Error("Verify Admin Request Validation failed",
					zap.String(constants.RequestIdLogKey, requestId),
				)
			}
			errResp := AuthModels.ErrorResponse{
				Message:   requestValidationFailedMessage,
				ErrorCode: 400,
				Errors:    validationErrors,
			}
			errorResponse, _ := json.Marshal(errResp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}

		if verifyAdminRequest.UserId == "" {
			if logger.IsErrorEnabled() {
				logger.Error("Verify Admin Request Validation failed",
					zap.String(constants.RequestIdLogKey, requestId),
				)
			}
			errResp := AuthModels.ErrorResponse{
				Message:   "Please provide valid user_id",
				ErrorCode: 400,
				Errors:    validationErrors,
			}
			errorResponse, _ := json.Marshal(errResp)
			writeErrorResponse(responseWriter, http.StatusBadRequest, errorResponse)
			return
		}
		ctx := context.WithValue(httpRequest.Context(), utils.RequestContextKeys.VerifyAdminRequestKey, verifyAdminRequest)
		next.ServeHTTP(responseWriter, httpRequest.WithContext(ctx))
	})
}
