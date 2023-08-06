package service

import (
	authModels "github.com/akgarg0472/urlshortener-auth-service/internal/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
)

var logger = Logger.GetLogger("authService.go")

func Login(requestId string, loginRequest authModels.LoginRequest) authModels.LoginResponse {
	logger.Info("[{}]: Processing Login Request -> {}", requestId, loginRequest)
	return authModels.LoginResponse{}
}

func Signup(requestId string, signupRequest authModels.SignupRequest) authModels.SignupResponse {
	logger.Info("[{}]: Processing Signup Request -> {}", requestId, signupRequest)
	return authModels.SignupResponse{}
}
