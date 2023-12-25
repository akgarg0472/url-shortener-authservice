package handler

import (
	"net/http"

	AuthService "github.com/akgarg0472/urlshortener-auth-service/internal/service/auth"
	AuthModels "github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
	Utils "github.com/akgarg0472/urlshortener-auth-service/utils"
)

var logger = Logger.GetLogger("authHandler.go")

// Handler Function to handle login request
func Login(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get("Request-ID")
	loginRequest := context.Value("loginRequest").(AuthModels.LoginRequest)

	logger.Trace("[{}]: Login request received on handler -> {}", requestId, loginRequest)

	loginResponse, loginError := AuthService.Login(requestId, loginRequest)

	sendResponseToClient(responseWriter, requestId, loginResponse, loginError, 200)
}

// Handler Function to handle signup request
func Signup(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get("Request-ID")
	signupRequest := context.Value("signupRequest").(AuthModels.SignupRequest)

	logger.Trace("[{}]: Signup request received on handler -> {}", requestId, signupRequest)

	signupResponse, signupError := AuthService.Signup(requestId, signupRequest)

	sendResponseToClient(responseWriter, requestId, signupResponse, signupError, 201)
}

// Handler Function to handle logout request
func Logout(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get("Request-ID")
	logoutRequest := context.Value("logoutRequest").(AuthModels.LogoutRequest)

	logger.Trace("[{}]: Logout request received on handler -> {}", requestId, logoutRequest)

	logoutResponse, logoutError := AuthService.Logout(requestId, logoutRequest)

	sendResponseToClient(responseWriter, requestId, logoutResponse, logoutError, 200)
}

// Handler Function to handle auth token validation request
func VerifyToken(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get("Request-ID")
	validateTokenRequest := context.Value("validateTokenRequest").(AuthModels.ValidateTokenRequest)

	logger.Trace("[{}]: Validate Token request received on handler -> {}", requestId, validateTokenRequest)

	validateTokenResponse, validateTokenError := AuthService.ValidateToken(requestId, validateTokenRequest)

	sendResponseToClient(responseWriter, requestId, validateTokenResponse, validateTokenError, 200)
}

// Handler Function to handle Forgot password request
func ForgotPassword(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get("Request-ID")
	forgotPasswordRequest := context.Value("forgotPasswordRequest").(AuthModels.ForgotPasswordRequest)

	logger.Trace("[{}]: Logout request received on handler -> {}", requestId, forgotPasswordRequest)

	forgotPasswordResponse, forgotPasswordError := AuthService.ForgotPassword(requestId, forgotPasswordRequest)

	sendResponseToClient(responseWriter, requestId, forgotPasswordResponse, forgotPasswordError, 200)
}

// Handler function to handle the verification of forgot password token verification check
func VerifyResetPassword(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	requestId := httpRequest.Header.Get("Request-ID")

	queryParams := httpRequest.URL.Query()

	logger.Trace("[{}]: Forgot Password verify request received on handler -> {}", requestId, queryParams)

	redirectUrl, err := AuthService.VerifyResetPassword(requestId, queryParams)

	if err != nil {
		sendResponseToClient(responseWriter, requestId, nil, err, 200)
		return
	}

	http.Redirect(responseWriter, httpRequest, redirectUrl, http.StatusSeeOther)
}

// Handler function to handle password reset (change) request
func ResetPassword(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get("Request-ID")
	resetPasswordRequest := context.Value("resetPasswordRequest").(AuthModels.ResetPasswordRequest)

	logger.Trace("[{}]: Reset Password request received on handler -> {}", requestId, resetPasswordRequest)

	resetPasswordResponse, resetPasswordError := AuthService.ResetPassword(requestId, resetPasswordRequest)

	sendResponseToClient(responseWriter, requestId, resetPasswordResponse, resetPasswordError, 200)
}

// Function to send response back to client
func sendResponseToClient(responseWriter http.ResponseWriter, requestId string, response interface{}, err *AuthModels.ErrorResponse, statusCode int) {
	if err != nil {
		errorJson, _ := utils.ConvertToJsonString(err)
		sendResponseWithStatusAndMessage(responseWriter, int(err.ErrorCode), errorJson)
		return
	}

	jsonResponse, jsonConvertError := Utils.ConvertToJsonString(response)

	if jsonConvertError != nil {
		logger.Error("[{}]: Error Converting Response to JSON: {}", requestId, jsonConvertError.Error())

		errorResponse := &AuthModels.ErrorResponse{
			Message:   "Internal Server Error",
			ErrorCode: 500,
		}

		errorResponseJson, _ := utils.ConvertToJsonString(errorResponse)
		sendResponseWithStatusAndMessage(responseWriter, http.StatusInternalServerError, errorResponseJson)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)
	responseWriter.Write([]byte(jsonResponse))
}

// Function to send error response to client with given status code and message
func sendResponseWithStatusAndMessage(responseWriter http.ResponseWriter, statusCode int, message string) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)
	responseWriter.Write([]byte(message))
}
