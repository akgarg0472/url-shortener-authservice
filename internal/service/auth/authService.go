package auth_service

import (
	"fmt"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	AuthDao "github.com/akgarg0472/urlshortener-auth-service/internal/dao"
	JwtService "github.com/akgarg0472/urlshortener-auth-service/internal/service/jwt"
	AuthModels "github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
)

var (
	logger = Logger.GetLogger("authService.go")
)

// Function to handle login request and generate jwt token
func Login(requestId string, loginRequest AuthModels.LoginRequest) (*AuthModels.LoginResponse, *AuthModels.ErrorResponse) {
	logger.Info("[{}]: Processing Login Request -> {}", requestId, loginRequest)

	user, err := AuthDao.GetUserByEmail(requestId, loginRequest.Email)

	if err != nil {
		logger.Error("[{}]: Error {} getting user by email -> {}", requestId, err.ErrorCode, err.Message)

		if err.ErrorCode == 404 {
			return nil, &AuthModels.ErrorResponse{
				Message:   "Invalid credentials",
				ErrorCode: 401,
			}
		}

		return nil, err
	}

	logger.Trace("[{}]: User -> {}", requestId, user)

	if !verifyPassword(loginRequest.Password, user.Password) {
		logger.Debug("[{}] invalid credentials provided", requestId)
		return nil, &AuthModels.ErrorResponse{Message: "Invalid credentials", ErrorCode: 401}
	}

	jwtToken, jwtError := JwtService.GetInstance().GenerateJwtToken(requestId, *user)

	if jwtError != nil {
		logger.Error("[{}]: Error while generating jwt token -> {}", requestId, jwtError)
		return nil, jwtError
	}

	return &AuthModels.LoginResponse{
		AccessToken: jwtToken,
		UserId:      user.Id,
		Name:        user.Email,
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

	logger.Debug("[{}]: User exists -> {}", requestId, strconv.FormatBool(userExists))

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

	user, saveError := AuthDao.SaveUser(requestId, signupRequest)

	if saveError != nil {
		logger.Error("[{}]: Error while saving user -> {}", requestId, saveError)
		return nil, saveError
	}

	if user == nil {
		logger.Error("[{}]: Something went wrong while saving user", requestId)
		return nil, utils.InternalServerErrorResponse()
	}

	return &AuthModels.SignupResponse{
		Message:    "User created successfully",
		StatusCode: 201,
	}, nil
}

// Function to handle logout request and invalidates the jwt token
func Logout(requestId string, logoutRequest AuthModels.LogoutRequest) (*AuthModels.LogoutResponse, *AuthModels.ErrorResponse) {
	logger.Info("[{}]: Processing Logout Request -> {}", requestId, logoutRequest)

	// todo: implement logic if required

	return &AuthModels.LogoutResponse{
		Message: "Logout successful",
	}, nil
}

// Function to handle validate token request and validates the jwt token
func ValidateToken(requestId string, validateTokenRequest AuthModels.ValidateTokenRequest) (*AuthModels.ValidateTokenResponse, *AuthModels.ErrorResponse) {
	logger.Debug("[{}]: Processing Validate Token Request -> {}", requestId, validateTokenRequest)

	token := validateTokenRequest.AuthToken
	userId := validateTokenRequest.UserId

	tokenValidateResp, err := JwtService.GetInstance().ValidateJwtToken(requestId, token, userId)

	if err != nil {
		logger.Error("[{}]: Error while validating jwt token -> {}", requestId, err)
		return nil, err
	}

	return tokenValidateResp, nil
}

// function to validate provided password against the encrypted password stored in DB
func verifyPassword(rawPassword string, encryptedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(rawPassword)) == nil
}
