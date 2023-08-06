package handler

import (
	"net/http"

	AuthService "github.com/akgarg0472/urlshortener-auth-service/internal/service/auth"
	AuthModels "github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	Utils "github.com/akgarg0472/urlshortener-auth-service/utils"
)

var logger = Logger.GetLogger("authHandler.go")

func Login(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get("Request-ID")
	loginRequest := context.Value("loginRequest").(AuthModels.LoginRequest)

	logger.Trace("[{}]: Login request received on handler -> {}", requestId, loginRequest)

	loginResponse := AuthService.Login(requestId, loginRequest)
	jsonResponse, err := Utils.ConvertToJsonString(loginResponse)

	if err != nil {
		logger.Error("[{}]: Error Converting Login Response to string: {}", requestId, err.Error())
		sendErrorResponse(responseWriter, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	responseWriter.Write([]byte(jsonResponse))
}

func Signup(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	context := httpRequest.Context()

	requestId := httpRequest.Header.Get("Request-ID")
	signupRequest := context.Value("signupRequest").(AuthModels.SignupRequest)

	logger.Trace("[{}]: Signup request received on handler -> {}", requestId, signupRequest)

	signupResponse := AuthService.Signup(requestId, signupRequest)
	jsonResponse, err := Utils.ConvertToJsonString(signupResponse)

	if err != nil {
		logger.Error("[{}]: Error while converting signup response to string -> {}", requestId, err)
		sendErrorResponse(responseWriter, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	responseWriter.Write([]byte(jsonResponse))
}

func sendErrorResponse(responseWriter http.ResponseWriter, statusCode int, message string) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)
	responseWriter.Write([]byte(message))
}
