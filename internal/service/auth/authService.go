package auth_service

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"

	authDao "github.com/akgarg0472/urlshortener-auth-service/internal/dao"
	notificationService "github.com/akgarg0472/urlshortener-auth-service/internal/service/notification"
	tokenService "github.com/akgarg0472/urlshortener-auth-service/internal/service/token"
	authModels "github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	utils "github.com/akgarg0472/urlshortener-auth-service/utils"
)

var (
	logger = Logger.GetLogger("authService.go")
)

// Function to handle login request and generate jwt token
func Login(requestId string, loginRequest authModels.LoginRequest) (*authModels.LoginResponse, *authModels.ErrorResponse) {
	logger.Info("[{}]: Processing Login Request -> {}", requestId, loginRequest)

	user, err := authDao.GetUserByEmail(requestId, loginRequest.Email)

	if err != nil {
		logger.Error("[{}]: Error {} getting user by email -> {}", requestId, err.ErrorCode, err.Message)

		if err.ErrorCode == 404 {
			return nil, &authModels.ErrorResponse{
				Message:   "Invalid credentials",
				ErrorCode: 401,
			}
		}

		return nil, err
	}

	logger.Trace("[{}]: User -> {}", requestId, user)

	if !verifyPassword(loginRequest.Password, user.Password) {
		logger.Error("[{}] invalid credentials provided", requestId)
		return nil, &authModels.ErrorResponse{Message: "Invalid credentials", ErrorCode: 401}
	}

	jwtToken, jwtError := tokenService.GetInstance().GenerateJwtToken(requestId, *user)

	if jwtError != nil {
		logger.Error("[{}]: Error generating auth token -> {}", requestId, jwtError)
		return nil, jwtError
	}

	authDao.UpdateTimestamp(requestId, loginRequest.Email, authDao.TIMESTAMP_TYPE_LAST_LOGIN_TIME)

	return &authModels.LoginResponse{
		AccessToken: jwtToken,
		UserId:      user.Id,
		Name:        user.Email,
	}, nil
}

// Function to handle signup request and save user in database
func Signup(requestId string, signupRequest authModels.SignupRequest) (*authModels.SignupResponse, *authModels.ErrorResponse) {
	logger.Info("[{}]: Processing Signup Request -> {}", requestId, signupRequest)

	userExists, userExistsError := authDao.CheckIfUserExistsByEmail(requestId, signupRequest.Email)

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

	user, saveError := authDao.SaveUser(requestId, signupRequest)

	if saveError != nil {
		logger.Error("[{}]: Error while saving user -> {}", requestId, saveError)
		return nil, saveError
	}

	if user == nil {
		logger.Error("[{}]: Something went wrong while saving user", requestId)
		return nil, utils.InternalServerErrorResponse()
	}

	notificationService.SendSignupSuccessEmail(requestId, user.Email, user.Name)

	return &authModels.SignupResponse{
		Message:    "Signup successful! You can now explore all of the exciting and amazing features",
		StatusCode: 201,
	}, nil
}

// Function to handle logout request and invalidates the jwt token
func Logout(requestId string, logoutRequest authModels.LogoutRequest) (*authModels.LogoutResponse, *authModels.ErrorResponse) {
	logger.Info("[{}]: Processing Logout Request -> {}", requestId, logoutRequest)

	// todo: implement logic if required

	return &authModels.LogoutResponse{
		Message: "Logout successful",
	}, nil
}

// Function to handle validate token request and validates the jwt token
func ValidateToken(requestId string, validateTokenRequest authModels.ValidateTokenRequest) (*authModels.ValidateTokenResponse, *authModels.ErrorResponse) {
	logger.Debug("[{}]: Processing Validate Token Request -> {}", requestId, validateTokenRequest)

	token := validateTokenRequest.AuthToken
	userId := validateTokenRequest.UserId

	tokenValidateResp, err := tokenService.GetInstance().ValidateJwtToken(requestId, token, userId)

	if err != nil {
		logger.Error("[{}]: Error while validating jwt token -> {}", requestId, err)
		return nil, err
	}

	return tokenValidateResp, nil
}

// Function to generate forgot password token and send forgot password email back to user
func GenerateAndSendForgotPasswordToken(requestId string, forgotPasswordRequest authModels.ForgotPasswordRequest) (*authModels.ForgotPasswordResponse, *authModels.ErrorResponse) {
	logger.Debug("[{}]: Processing forgot password Request -> {}", requestId, forgotPasswordRequest)

	email := forgotPasswordRequest.Email

	user, err := authDao.GetUserByEmail(requestId, email)

	if err != nil {
		if err.ErrorCode == 404 {
			err.ErrorCode = 400
		}

		return nil, &authModels.ErrorResponse{
			Message:   err.Message,
			Errors:    err.Errors,
			ErrorCode: err.ErrorCode,
		}
	}

	forgotPasswordToken, err := tokenService.GetInstance().GenerateForgotPasswordToken(requestId, *user)

	if err != nil {
		return nil, err
	}

	// store token in database for corresponding user
	dbUpdated, dbUpdateError := authDao.UpdateForgotPasswordToken(requestId, user.Email, forgotPasswordToken)

	if dbUpdateError != nil {
		return nil, dbUpdateError
	}

	if !dbUpdated {
		return nil, utils.InternalServerErrorResponse()
	}

	// now generate forgot password link which will be sent on user's email
	tokenResetLink := utils.GenerateForgotPasswordLink(user.Email, forgotPasswordToken)

	// send email to user and return success response
	notificationService.SendForgotPasswordEmail(requestId, user.Email, user.Name, tokenResetLink)

	return &authModels.ForgotPasswordResponse{
		Success:    true,
		Message:    "We have sent an email to " + email + " with steps to reset your password. Please follow email to continue",
		StatusCode: 200,
	}, nil
}

// Function to validate forgot password token and return redirect URL to reset password UI page
func VerifyResetPasswordToken(requestId string, queryParams url.Values) (string, *authModels.ErrorResponse) {
	emailParam := queryParams["email"]
	tokenParam := queryParams["token"]

	resetPasswordValidationError := utils.ValidateResetPasswordRequestQueryParams(emailParam, tokenParam)

	if resetPasswordValidationError != nil {
		return "", resetPasswordValidationError
	}

	email := emailParam[0]
	token := tokenParam[0]

	tokenValidationError := tokenService.GetInstance().ValidateForgotPasswordToken(requestId, token)

	if tokenValidationError != nil {
		return "", tokenValidationError
	}

	forgotPasswordTokenFromDatabase, fptfdError := authDao.GetForgotPasswordToken(requestId, email)

	if fptfdError != nil {
		return "", fptfdError
	}

	if forgotPasswordTokenFromDatabase != token {
		logger.Error("[{}] Forgot Token not found in Database", requestId)

		return "", &authModels.ErrorResponse{
			Message:   "Invalid forgot password token. Please try again",
			ErrorCode: 400,
		}
	}

	redirectUrl := utils.GenerateForgotPasswordTokenRedirectUrl(email, token)

	logger.Debug("[{}] Redirect URL generated is: {}", requestId, redirectUrl)

	return redirectUrl, nil
}

// Function to actually reset password from forgot-password UI page
func ResetPassword(requestId string, resetPasswordRequest authModels.ResetPasswordRequest) (*authModels.ResetPasswordResponse, *authModels.ErrorResponse) {
	logger.Info("[{}]: Processing Reset password Request", requestId)

	email := resetPasswordRequest.Email
	resetPasswordToken := resetPasswordRequest.ResetPasswordToken
	password := resetPasswordRequest.Password
	confirmPassword := resetPasswordRequest.ConfirmPassword

	// verify and match password
	if strings.TrimSpace(password) != strings.TrimSpace(confirmPassword) {
		logger.Error("[{}] Password & confirm Passwords mismatch", requestId)

		return nil, &authModels.ErrorResponse{
			Message:   "Password & confirm Passwords mismatch",
			ErrorCode: 400,
		}
	}

	// fetch forgot password token from DB
	forgotPasswordTokenFromDatabase, fptfdError := authDao.GetForgotPasswordToken(requestId, email)

	if fptfdError != nil {
		logger.Error("[{}] error fetching forgot password token from DB", requestId)
		return nil, fptfdError
	}

	// match provided token with DB token again for double check
	if strings.TrimSpace(forgotPasswordTokenFromDatabase) != strings.TrimSpace(resetPasswordToken) {
		logger.Error("[{}] forgot token from DB doesn't match with token provided", requestId)

		return nil, &authModels.ErrorResponse{
			Message:   "Invalid token provided",
			ErrorCode: 400,
		}
	}

	// reset password
	hashedPassword, bcryptError := bcrypt.GenerateFromPassword([]byte(password), 14)

	if bcryptError != nil {
		logger.Error("[{}]: Error while hashing password -> {}", requestId, bcryptError)
		return nil, utils.InternalServerErrorResponse()
	}

	isPasswordUpdated, passwordUpdateErr := authDao.UpdatePassword(requestId, email, string(hashedPassword))

	if passwordUpdateErr != nil {
		return nil, passwordUpdateErr
	}

	if !isPasswordUpdated {
		return nil, utils.InternalServerErrorResponse()
	}

	notificationService.SendPasswordChangeSuccessEmail(requestId, email)

	return &authModels.ResetPasswordResponse{
		Success:    true,
		Message:    "Password changed successfully",
		StatusCode: 200,
	}, nil
}
