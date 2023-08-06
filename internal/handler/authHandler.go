package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	authModels "github.com/akgarg0472/urlshortener-auth-service/internal/model"
	authService "github.com/akgarg0472/urlshortener-auth-service/internal/service/auth"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
)

var logger = Logger.GetLogger("authHandler.go")

func Login(w http.ResponseWriter, r *http.Request) {
	context := r.Context()

	requestId := r.Header.Get("Request-ID")

	loginRequest := context.Value("loginRequest").(authModels.LoginRequest)
	fmt.Println("Request id: " + requestId)

	logger.Trace("[{}]: Login request received on handler -> {}", requestId, loginRequest)

	loginResponse := authService.Login(requestId, loginRequest)

	jsonResponse, err := stringify(loginResponse)

	if err != nil {
		logger.Error("")
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Write([]byte(jsonResponse))
}

func Signup(w http.ResponseWriter, r *http.Request) {
	context := r.Context()

	requestId := r.Header.Get("Request-ID")
	signupRequest := context.Value("signupRequest").(authModels.SignupRequest)

	logger.Trace("[{}]: Signup request received on handler -> {}", requestId, signupRequest)

	signupResponse := authService.Signup(requestId, signupRequest)

	jsonResponse, err := stringify(signupResponse)

	if err != nil {
		logger.Error("[{}]: Error while converting signup response to string -> {}", requestId, err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Write([]byte(jsonResponse))
}

func stringify(response interface{}) (string, error) {
	responseBytes, err := json.Marshal(response)
	return string(responseBytes), err
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
