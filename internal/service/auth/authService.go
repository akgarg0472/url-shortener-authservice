package auth_service

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"

	AuthDao "github.com/akgarg0472/urlshortener-auth-service/internal/dao"
	TokenService "github.com/akgarg0472/urlshortener-auth-service/internal/service/token"
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

	jwtToken, jwtError := TokenService.GetInstance().GenerateJwtToken(requestId, *user)

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

	tokenValidateResp, err := TokenService.GetInstance().ValidateJwtToken(requestId, token, userId)

	if err != nil {
		logger.Error("[{}]: Error while validating jwt token -> {}", requestId, err)
		return nil, err
	}

	return tokenValidateResp, nil
}

// Function to generate forgot password token and send forgot password email back to user
func ForgotPassword(requestId string, forgotPasswordRequest AuthModels.ForgotPasswordRequest) (*AuthModels.ForgotPasswordResponse, *AuthModels.ErrorResponse) {
	logger.Debug("[{}]: Processing forgot password Request -> {}", requestId, forgotPasswordRequest)

	email := forgotPasswordRequest.Email

	user, err := AuthDao.GetUserByEmail(requestId, email)

	if err != nil {
		if err.ErrorCode == 404 {
			err.ErrorCode = 400
		}

		return nil, &AuthModels.ErrorResponse{
			Message:   err.Message,
			Errors:    err.Errors,
			ErrorCode: err.ErrorCode,
		}
	}

	forgotPasswordToken, err := TokenService.GetInstance().GenerateForgotPasswordToken(requestId, *user)

	if err != nil {
		return nil, err
	}

	// store token in database for corresponding user
	dbUpdated, dbUpdateError := AuthDao.UpdateForgotPasswordToken(requestId, user.Id, forgotPasswordToken)

	if dbUpdateError != nil {
		return nil, dbUpdateError
	}

	if !dbUpdated {
		return nil, utils.InternalServerErrorResponse()
	}

	// now generate reset password link which will be sent on user's email
	tokenResetLink := utils.GenerateResetPasswordLink(user.Email, forgotPasswordToken)

	// send email to user and return success response
	emailSent := utils.SendForgotPasswordEmailToUser(user.FirstName+user.LastName, user.Email, tokenResetLink)

	if !emailSent {
		return nil, utils.InternalServerErrorResponse()
	}

	return &AuthModels.ForgotPasswordResponse{
		Success:    true,
		Message:    "We have sent an email to " + email + " with steps to reset your password",
		StatusCode: 200,
	}, nil
}

// Function to validate forgot password token and redirect to reset password UI page
func VerifyResetPassword(requestId string, queryParams url.Values) (string, *AuthModels.ErrorResponse) {
	emailParam := queryParams["email"]
	tokenParam := queryParams["token"]

	resetPasswordValidationError := utils.ValidateResetPasswordRequestQueryParams(emailParam, tokenParam)

	if resetPasswordValidationError != nil {
		return "", resetPasswordValidationError
	}

	email := emailParam[0]
	token := tokenParam[0]

	tokenValidationError := TokenService.GetInstance().ValidateForgotPasswordToken(requestId, token)

	if tokenValidationError != nil {
		return "", tokenValidationError
	}

	forgotPasswordTokenFromDatabase, fptfdError := AuthDao.GetForgotPasswordToken(requestId, email)

	if fptfdError != nil {
		return "", fptfdError
	}

	if forgotPasswordTokenFromDatabase != token {
		logger.Error("[{}] Forgot Token not found in Database", requestId)

		return "", &AuthModels.ErrorResponse{
			Message:   "Invalid forgot password token. Please try again by requesting for reset password",
			ErrorCode: 400,
		}
	}

	redirectUrl := utils.GenerateForgotPasswordTokenRedirectUrl(email, token)

	logger.Debug("[{}] Redirect URL generated is: {}", requestId, redirectUrl)

	return redirectUrl, nil
}

// Function to actually reset password
func ResetPassword(requestId string, resetPasswordRequest AuthModels.ResetPasswordRequest) (*AuthModels.ResetPasswordResponse, *AuthModels.ErrorResponse) {
	logger.Debug("[{}]: Processing Reset password Request", requestId)

	email := resetPasswordRequest.Email
	resetPasswordToken := resetPasswordRequest.ResetPasswordToken
	password := resetPasswordRequest.Password
	confirmPassword := resetPasswordRequest.ConfirmPassword

	// verify and match password
	if strings.TrimSpace(password) != strings.TrimSpace(confirmPassword) {
		logger.Error("[{}] Password & confirm Passwords mismatch", requestId)

		return nil, &AuthModels.ErrorResponse{
			Message:   "Password & confirm Passwords mismatch",
			ErrorCode: 400,
		}
	}

	// fetch forgot password token from DB
	forgotPasswordTokenFromDatabase, fptfdError := AuthDao.GetForgotPasswordToken(requestId, email)

	if fptfdError != nil {
		logger.Error("[{}] error fetching forgot password token from DB", requestId)
		return nil, fptfdError
	}

	// match provided token with DB token again for double check
	if strings.TrimSpace(forgotPasswordTokenFromDatabase) != strings.TrimSpace(resetPasswordToken) {
		logger.Error("[{}] forgot token from DB doesn't match with token provided", requestId)

		return nil, &AuthModels.ErrorResponse{
			Message:   "Invalid token provided",
			ErrorCode: 400,
		}
	}

	// reset passwords
	hashedPassword, bcryptError := bcrypt.GenerateFromPassword([]byte(password), 14)

	if bcryptError != nil {
		logger.Error("[{}]: Error while hashing password -> {}", requestId, bcryptError)
		return nil, utils.InternalServerErrorResponse()
	}

	isPasswordUpdated, passwordUpdateErr := AuthDao.UpdatePassword(requestId, email, string(hashedPassword))

	if passwordUpdateErr != nil {
		return nil, passwordUpdateErr
	}

	if !isPasswordUpdated {
		return nil, &AuthModels.ErrorResponse{
			Message:   "Error resetting password",
			ErrorCode: 400,
		}
	}

	// reset token to default empty string
	AuthDao.UpdateForgotPasswordToken(requestId, email, "")

	return &AuthModels.ResetPasswordResponse{
		Success:    true,
		Message:    "Password reset successfully",
		StatusCode: 200,
	}, nil
}

// function to validate provided password against the encrypted password stored in DB
func verifyPassword(rawPassword string, encryptedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(rawPassword)) == nil
}
