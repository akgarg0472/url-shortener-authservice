package handler

import (
	"net/http"
	"time"

	"github.com/akgarg0472/urlshortener-auth-service/constants"
	"github.com/akgarg0472/urlshortener-auth-service/internal/logger"
	auth_service "github.com/akgarg0472/urlshortener-auth-service/internal/service/auth"
	"github.com/akgarg0472/urlshortener-auth-service/model"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
	"go.uber.org/zap"
)

// LoginHandler Handler Function to handle login request
func LoginHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get(constants.RequestIdHeaderName)
	loginRequest := context.Value(utils.RequestContextKeys.LoginRequestKey).(model.LoginRequest)

	if logger.IsDebugEnabled() {
		logger.Debug("User login attempt received via email/password",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.Any(constants.RequestLogKey, loginRequest),
		)
	}

	loginResponse, loginError := auth_service.LoginWithEmailPassword(requestId, loginRequest)

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    loginResponse.AccessToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(responseWriter, cookie)

	sendResponseToClient(responseWriter, requestId, loginResponse, loginError, 200)
}

// SignupHandler Handler Function to handle signup request
func SignupHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get(constants.RequestIdHeaderName)
	signupRequest := context.Value(utils.RequestContextKeys.SignupRequestKey).(model.SignupRequest)

	if logger.IsDebugEnabled() {
		logger.Debug("User signup request received",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.Any(constants.RequestLogKey, signupRequest),
		)
	}

	signupResponse, signupError := auth_service.Signup(requestId, signupRequest)

	sendResponseToClient(responseWriter, requestId, signupResponse, signupError, 201)
}

// LogoutHandler Handler Function to handle logout request
func LogoutHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get(constants.RequestIdHeaderName)
	logoutRequest := context.Value(utils.RequestContextKeys.LogoutRequestKey).(model.LogoutRequest)

	if logger.IsDebugEnabled() {
		logger.Debug("User logout request received",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.Any(constants.RequestLogKey, logoutRequest),
		)
	}

	logoutResponse, logoutError := auth_service.Logout(requestId, logoutRequest)

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(-1 * time.Hour),
	}

	http.SetCookie(responseWriter, cookie)

	sendResponseToClient(responseWriter, requestId, logoutResponse, logoutError, 200)
}

// VerifyTokenHandler Handler Function to handle auth token validation request
func VerifyTokenHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get(constants.RequestIdHeaderName)
	validateTokenRequest := context.Value(utils.RequestContextKeys.ValidateTokenRequestKey).(model.ValidateTokenRequest)

	if logger.IsDebugEnabled() {
		logger.Debug("Token validation request received",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.Any(constants.RequestLogKey, validateTokenRequest),
		)
	}

	validateTokenResponse, validateTokenError := auth_service.ValidateToken(requestId, validateTokenRequest)

	sendResponseToClient(responseWriter, requestId, validateTokenResponse, validateTokenError, 200)
}

// ForgotPasswordHandler Handler Function to handle Forgot password request
func ForgotPasswordHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get(constants.RequestIdHeaderName)
	forgotPasswordRequest := context.Value(utils.RequestContextKeys.ForgotPasswordRequestKey).(model.ForgotPasswordRequest)

	if logger.IsDebugEnabled() {
		logger.Debug("User logout request received",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.Any(constants.RequestLogKey, forgotPasswordRequest),
		)
	}

	forgotPasswordResponse, forgotPasswordError := auth_service.GenerateAndSendForgotPasswordToken(requestId, forgotPasswordRequest)

	sendResponseToClient(responseWriter, requestId, forgotPasswordResponse, forgotPasswordError, 200)
}

// VerifyResetPasswordHandler Handler function to handle the verification of forgot password token verification check
func VerifyResetPasswordHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	requestId := httpRequest.Header.Get(constants.RequestIdHeaderName)

	queryParams := httpRequest.URL.Query()

	if logger.IsDebugEnabled() {
		logger.Debug("Forgot password verification request received",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.Any(constants.RequestLogKey, queryParams),
		)
	}

	redirectUrl, err := auth_service.VerifyResetPasswordToken(requestId, queryParams)

	if err != nil {
		sendResponseToClient(responseWriter, requestId, nil, err, 200)
		return
	}

	http.Redirect(responseWriter, httpRequest, redirectUrl, http.StatusSeeOther)
}

// ResetPasswordHandler Handler function to handle password reset (change) request
func ResetPasswordHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get(constants.RequestIdHeaderName)
	resetPasswordRequest := context.Value(utils.RequestContextKeys.ResetPasswordRequestKey).(model.ResetPasswordRequest)

	if logger.IsDebugEnabled() {
		logger.Debug("Password reset request received",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.Any(constants.RequestLogKey, resetPasswordRequest),
		)
	}

	resetPasswordResponse, resetPasswordError := auth_service.ResetPassword(requestId, resetPasswordRequest)

	sendResponseToClient(responseWriter, requestId, resetPasswordResponse, resetPasswordError, 200)
}

// VerifyAdminHandler Handler function to handle verify admin request
func VerifyAdminHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get(constants.RequestIdHeaderName)
	verifyAdminRequest := context.Value(utils.RequestContextKeys.VerifyAdminRequestKey).(model.VerifyAdminRequest)

	if logger.IsDebugEnabled() {
		logger.Debug("Admin verification request received",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.Any(constants.RequestLogKey, verifyAdminRequest),
		)
	}

	verifyAdminResponse, verifyAdminError := auth_service.VerifyAdmin(requestId, verifyAdminRequest)

	sendResponseToClient(responseWriter, requestId, verifyAdminResponse, verifyAdminError, 200)
}
