package service

import (
	AuthModels "github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
)

var logger = Logger.GetLogger("authService.go")

func Login(requestId string, loginRequest AuthModels.LoginRequest) AuthModels.LoginResponse {
	logger.Info("[{}]: Processing Login Request -> {}", requestId, loginRequest)
	return AuthModels.LoginResponse{}
}

func Signup(requestId string, signupRequest AuthModels.SignupRequest) AuthModels.SignupResponse {
	logger.Info("[{}]: Processing Signup Request -> {}", requestId, signupRequest)
	return AuthModels.SignupResponse{}
}
