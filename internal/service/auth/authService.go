package auth_service

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"

	AuthDao "github.com/akgarg0472/urlshortener-auth-service/internal/dao"
	JwtService "github.com/akgarg0472/urlshortener-auth-service/internal/service/jwt"
	AuthModels "github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
	"github.com/google/uuid"
)

var (
	logger = Logger.GetLogger("authService.go")
)

// Function to handle login request and generate jwt token
func Login(requestId string, loginRequest AuthModels.LoginRequest) (*AuthModels.LoginResponse, *AuthModels.ErrorResponse) {
	logger.Info("[{}]: Processing Login Request -> {}", requestId, loginRequest)

	user, err := AuthDao.GetUserByEmail(requestId, loginRequest.Email)

	if err != nil {
		logger.Error("[{}]: Error while getting user by email -> {}", requestId, err)
		return nil, err
	}

	logger.Debug("[{}]: User -> {}", requestId, user)

	jwtToken, jwtError := JwtService.GetInstance().GenerateJwtToken(requestId, *user)

	if jwtError != nil {
		logger.Error("[{}]: Error while generating jwt token -> {}", requestId, jwtError)
		return nil, jwtError
	}

	return &AuthModels.LoginResponse{
		AccessToken: jwtToken,
	}, nil
}

// Function to handle signup request and save user in database
func Signup(requestId string, signupRequest AuthModels.SignupRequest) (*AuthModels.SignupResponse, *AuthModels.ErrorResponse) {
	logger.Info("[{}]: Processing Signup Request -> {}", requestId, signupRequest)

	userExists, userExistsError := AuthDao.CheckIfUserExistsByEmail(requestId, signupRequest.Email)

	if userExistsError != nil {
		logger.Error("[{}]: Error while checking user exists -> {}", requestId, userExistsError)
		return nil, userExistsError
	}

	logger.Trace("[{}]: User exists -> {}", requestId, strconv.FormatBool(userExists))

	if userExists {
		logger.Error("[{}]: Email already exists -> {}", requestId, signupRequest.Email)
		return nil, utils.GetErrorResponse(fmt.Sprintf("User already exists with email: %s", signupRequest.Email), 409)
	}

	hashedPassword, bcryptError := bcrypt.GenerateFromPassword([]byte(signupRequest.Password), 14)

	if bcryptError != nil {
		logger.Error("[{}]: Error while hashing password -> {}", requestId, bcryptError)
		return nil, utils.InternalServerErrorResponse()
	}

	signupRequest.Password = string(hashedPassword)
	signupRequest.UserId = strings.ReplaceAll(uuid.New().String(), "-", "")

	saveSuccess, saveError := AuthDao.SaveUser(requestId, signupRequest)

	if saveError != nil {
		logger.Error("[{}]: Error while saving user -> {}", requestId, saveError)
		return nil, saveError
	}

	if !saveSuccess {
		logger.Error("[{}]: Something went wrong while saving user", requestId)
		return nil, utils.InternalServerErrorResponse()
	}

	return &AuthModels.SignupResponse{
		Message: "User created successfully",
		UserId:  signupRequest.UserId,
	}, nil
}

// Function to handle logout request and invalidates the jwt token
func Logout(requestId string, logoutRequest AuthModels.LogoutRequest) (*AuthModels.LogoutResponse, *AuthModels.ErrorResponse) {
	logger.Info("[{}]: Processing Logout Request -> {}", requestId, logoutRequest)

	token := logoutRequest.AuthToken
	userId := logoutRequest.UserId

	err := JwtService.GetInstance().InvalidateJwtToken(requestId, token, userId)

	if err != nil {
		logger.Error("[{}]: Error while invalidating jwt token -> {}", requestId, err)
		return nil, err
	}

	return &AuthModels.LogoutResponse{
		Message: "Logout successful",
	}, nil
}

// Function to handle validate token request and validates the jwt token
func ValidateToken(requestId string, validateTokenRequest AuthModels.ValidateTokenRequest) (*AuthModels.ValidateTokenResponse, *AuthModels.ErrorResponse) {
	logger.Debug("[{}]: Processing Validate Token Request -> {}", requestId, validateTokenRequest)

	token := validateTokenRequest.AuthToken
	userId := validateTokenRequest.UserId

	err := JwtService.GetInstance().ValidateJwtToken(requestId, token, userId)

	if err != nil {
		logger.Error("[{}]: Error while validating jwt token -> {}", requestId, err)
		return nil, err
	}

	return &AuthModels.ValidateTokenResponse{
		Message: "Token validated successfully",
	}, nil
}
