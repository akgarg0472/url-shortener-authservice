package handler

import (
	"net/http"

	authService "github.com/akgarg0472/urlshortener-auth-service/internal/service/auth"
	authModels "github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
)

var authLogger = Logger.GetLogger("authHandler.go")

// Handler Function to handle login request
func LoginHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get("Request-ID")
	loginRequest := context.Value("loginRequest").(authModels.LoginRequest)

	authLogger.Trace("[{}]: Login request received on handler -> {}", requestId, loginRequest)

	loginResponse, loginError := authService.Login(requestId, loginRequest)

	sendResponseToClient(responseWriter, requestId, loginResponse, loginError, 200)
}

// Handler Function to handle signup request
func SignupHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get("Request-ID")
	signupRequest := context.Value("signupRequest").(authModels.SignupRequest)

	authLogger.Trace("[{}]: Signup request received on handler -> {}", requestId, signupRequest)

	signupResponse, signupError := authService.Signup(requestId, signupRequest)

	sendResponseToClient(responseWriter, requestId, signupResponse, signupError, 201)
}

// Handler Function to handle logout request
func LogoutHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get("Request-ID")
	logoutRequest := context.Value("logoutRequest").(authModels.LogoutRequest)

	authLogger.Trace("[{}]: Logout request received on handler -> {}", requestId, logoutRequest)

	logoutResponse, logoutError := authService.Logout(requestId, logoutRequest)

	sendResponseToClient(responseWriter, requestId, logoutResponse, logoutError, 200)
}

// Handler Function to handle auth token validation request
func VerifyTokenHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get("Request-ID")
	validateTokenRequest := context.Value("validateTokenRequest").(authModels.ValidateTokenRequest)

	authLogger.Trace("[{}]: Validate Token request received on handler -> {}", requestId, validateTokenRequest)

	validateTokenResponse, validateTokenError := authService.ValidateToken(requestId, validateTokenRequest)

	sendResponseToClient(responseWriter, requestId, validateTokenResponse, validateTokenError, 200)
}

// Handler Function to handle Forgot password request
func ForgotPasswordHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get("Request-ID")
	forgotPasswordRequest := context.Value("forgotPasswordRequest").(authModels.ForgotPasswordRequest)

	authLogger.Trace("[{}]: Logout request received on handler -> {}", requestId, forgotPasswordRequest)

	forgotPasswordResponse, forgotPasswordError := authService.GenerateAndSendForgotPasswordToken(requestId, forgotPasswordRequest)

	sendResponseToClient(responseWriter, requestId, forgotPasswordResponse, forgotPasswordError, 200)
}

// Handler function to handle the verification of forgot password token verification check
func VerifyResetPasswordHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	requestId := httpRequest.Header.Get("Request-ID")

	queryParams := httpRequest.URL.Query()

	authLogger.Trace("[{}]: Forgot Password verify request received on handler -> {}", requestId, queryParams)

	redirectUrl, err := authService.VerifyResetPasswordToken(requestId, queryParams)

	if err != nil {
		sendResponseToClient(responseWriter, requestId, nil, err, 200)
		return
	}

	http.Redirect(responseWriter, httpRequest, redirectUrl, http.StatusSeeOther)
}

// Handler function to handle password reset (change) request
func ResetPasswordHandler(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get("Request-ID")
	resetPasswordRequest := context.Value("resetPasswordRequest").(authModels.ResetPasswordRequest)

	authLogger.Trace("[{}]: Reset Password request received on handler -> {}", requestId, resetPasswordRequest)

	resetPasswordResponse, resetPasswordError := authService.ResetPassword(requestId, resetPasswordRequest)

	sendResponseToClient(responseWriter, requestId, resetPasswordResponse, resetPasswordError, 200)
}
