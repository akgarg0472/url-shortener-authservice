package handler

import (
	"net/http"

	authService "github.com/akgarg0472/urlshortener-auth-service/internal/service/auth"
	authModels "github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
)

var authLogger = Logger.GetLogger("authHandler.go")
var requestIdHeader = "X-Request-Id"

// LoginHandler Handler Function to handle login request
func LoginHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get(requestIdHeader)
	loginRequest := context.Value(utils.RequestContextKeys.LoginRequestKey).(authModels.LoginRequest)

	authLogger.Trace("[{}]: LoginWithEmailPassword request received on handler -> {}", requestId, loginRequest)

	loginResponse, loginError := authService.LoginWithEmailPassword(requestId, loginRequest)

	sendResponseToClient(responseWriter, requestId, loginResponse, loginError, 200)
}

// SignupHandler Handler Function to handle signup request
func SignupHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get(requestIdHeader)
	signupRequest := context.Value(utils.RequestContextKeys.SignupRequestKey).(authModels.SignupRequest)

	authLogger.Trace("[{}]: Signup request received on handler -> {}", requestId, signupRequest)

	signupResponse, signupError := authService.Signup(requestId, signupRequest)

	sendResponseToClient(responseWriter, requestId, signupResponse, signupError, 201)
}

// LogoutHandler Handler Function to handle logout request
func LogoutHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get(requestIdHeader)
	logoutRequest := context.Value(utils.RequestContextKeys.LogoutRequestKey).(authModels.LogoutRequest)

	authLogger.Trace("[{}]: Logout request received on handler -> {}", requestId, logoutRequest)

	logoutResponse, logoutError := authService.Logout(requestId, logoutRequest)

	sendResponseToClient(responseWriter, requestId, logoutResponse, logoutError, 200)
}

// VerifyTokenHandler Handler Function to handle auth token validation request
func VerifyTokenHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get(requestIdHeader)
	validateTokenRequest := context.Value(utils.RequestContextKeys.ValidateTokenRequestKey).(authModels.ValidateTokenRequest)

	authLogger.Trace("[{}]: Validate Token request received on handler -> {}", requestId, validateTokenRequest)

	validateTokenResponse, validateTokenError := authService.ValidateToken(requestId, validateTokenRequest)

	sendResponseToClient(responseWriter, requestId, validateTokenResponse, validateTokenError, 200)
}

// ForgotPasswordHandler Handler Function to handle Forgot password request
func ForgotPasswordHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get(requestIdHeader)
	forgotPasswordRequest := context.Value(utils.RequestContextKeys.ForgotPasswordRequestKey).(authModels.ForgotPasswordRequest)

	authLogger.Trace("[{}]: Logout request received on handler -> {}", requestId, forgotPasswordRequest)

	forgotPasswordResponse, forgotPasswordError := authService.GenerateAndSendForgotPasswordToken(requestId, forgotPasswordRequest)

	sendResponseToClient(responseWriter, requestId, forgotPasswordResponse, forgotPasswordError, 200)
}

// VerifyResetPasswordHandler Handler function to handle the verification of forgot password token verification check
func VerifyResetPasswordHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	requestId := httpRequest.Header.Get(requestIdHeader)

	queryParams := httpRequest.URL.Query()

	authLogger.Trace("[{}]: Forgot Password verify request received on handler -> {}", requestId, queryParams)

	redirectUrl, err := authService.VerifyResetPasswordToken(requestId, queryParams)

	if err != nil {
		sendResponseToClient(responseWriter, requestId, nil, err, 200)
		return
	}

	http.Redirect(responseWriter, httpRequest, redirectUrl, http.StatusSeeOther)
}

// ResetPasswordHandler Handler function to handle password reset (change) request
func ResetPasswordHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get(requestIdHeader)
	resetPasswordRequest := context.Value(utils.RequestContextKeys.ResetPasswordRequestKey).(authModels.ResetPasswordRequest)

	authLogger.Trace("[{}]: Reset Password request received on handler -> {}", requestId, resetPasswordRequest)

	resetPasswordResponse, resetPasswordError := authService.ResetPassword(requestId, resetPasswordRequest)

	sendResponseToClient(responseWriter, requestId, resetPasswordResponse, resetPasswordError, 200)
}

// VerifyAdminHandler Handler function to handle verify admin request
func VerifyAdminHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get(requestIdHeader)
	verifyAdminRequest := context.Value(utils.RequestContextKeys.VerifyAdminRequestKey).(authModels.VerifyAdminRequest)

	authLogger.Trace("[{}]: Verify Admin request received on handler -> {}", requestId, verifyAdminRequest)

	verifyAdminResponse, verifyAdminError := authService.VerifyAdmin(requestId, verifyAdminRequest)

	sendResponseToClient(responseWriter, requestId, verifyAdminResponse, verifyAdminError, 200)
}
