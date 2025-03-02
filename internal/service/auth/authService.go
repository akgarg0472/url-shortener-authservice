package auth_service

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/akgarg0472/urlshortener-auth-service/constants"
	"github.com/akgarg0472/urlshortener-auth-service/model"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"golang.org/x/crypto/bcrypt"

	authDao "github.com/akgarg0472/urlshortener-auth-service/internal/dao/auth"
	"github.com/akgarg0472/urlshortener-auth-service/internal/entity"
	"github.com/akgarg0472/urlshortener-auth-service/internal/logger"
	kafka_service "github.com/akgarg0472/urlshortener-auth-service/internal/service/kafka"
	notificationService "github.com/akgarg0472/urlshortener-auth-service/internal/service/notification"
	tokenService "github.com/akgarg0472/urlshortener-auth-service/internal/service/token"
	authModels "github.com/akgarg0472/urlshortener-auth-service/model"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
)

// LoginWithEmailPassword Function to handle login request using email & password and generate JWT token
func LoginWithEmailPassword(requestId string, loginRequest authModels.LoginRequest) (*authModels.LoginResponse, *authModels.ErrorResponse) {
	if logger.IsInfoEnabled() {
		logger.Info(
			"Processing LoginWithEmailPassword Request",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.Any("loginRequest", loginRequest),
		)
	}

	user, err := authDao.GetUserByEmail(requestId, loginRequest.Email)

	if err != nil {
		if logger.IsErrorEnabled() {
			logger.Error(
				"Error getting user by email",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Any(constants.ErrorCodeLogKey, err.ErrorCode),
				zap.Any(constants.ErrorMessageLogKey, err.Message),
			)
		}

		if err.ErrorCode == 404 {
			return nil, &authModels.ErrorResponse{
				Message:   "Invalid credentials",
				ErrorCode: 401,
			}
		}

		return nil, err
	}

	if logger.IsDebugEnabled() {
		logger.Debug(
			"User details",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.Any("user", user),
		)
	}

	if user.LoginType != constants.UserEntityLoginTypeEmailAndPassword {
		if logger.IsInfoEnabled() {
			logger.Info(
				"User is not registered using email and password",
				zap.String(constants.RequestIdLogKey, requestId),
			)
		}
		return nil, &authModels.ErrorResponse{
			Message:   fmt.Sprintf("Your account is registered using %s OAuth and does not have a password. Please log in using %s OAuth.", user.OAuthProvider, user.OAuthProvider),
			ErrorCode: 400,
		}
	}

	if !verifyPassword(loginRequest.Password, user.Password) {
		if logger.IsErrorEnabled() {
			logger.Error(
				"Invalid credentials provided",
				zap.String(constants.RequestIdLogKey, requestId),
			)
		}
		return nil, &authModels.ErrorResponse{Message: "Invalid credentials", ErrorCode: 401}
	}

	jwtToken, jwtError := tokenService.GetInstance().GenerateJwtToken(requestId, *user)

	if jwtError != nil {
		if logger.IsErrorEnabled() {
			logger.Error(
				"Error generating auth token",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Int16(constants.ErrorCodeLogKey, jwtError.ErrorCode),
				zap.Any(constants.ErrorMessageLogKey, jwtError.Message),
			)
		}
		return nil, jwtError
	}

	authDao.UpdateTimestamp(requestId, loginRequest.Email, authDao.TimestampTypeLastLoginTime)

	return &authModels.LoginResponse{
		AccessToken: jwtToken,
		UserId:      user.Id,
		Name:        user.Name,
		Email:       user.Email,
		LoginType:   string(user.LoginType),
	}, nil
}

// Signup Function to handle signup request and save user in database
func Signup(requestId string, signupRequest authModels.SignupRequest) (*authModels.SignupResponse, *authModels.ErrorResponse) {
	if logger.IsInfoEnabled() {
		logger.Info(
			"Processing Signup Request",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.Any("signupRequest", signupRequest),
		)
	}

	userExists, userExistsError := authDao.CheckIfUserExistsByEmail(requestId, signupRequest.Email)

	if userExistsError != nil {
		if logger.IsErrorEnabled() {
			logger.Error(
				"Error while checking user exists",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Int16(constants.ErrorCodeLogKey, userExistsError.ErrorCode),
				zap.Any(constants.ErrorMessageLogKey, userExistsError.Message),
			)
		}
		return nil, userExistsError
	}

	if logger.IsInfoEnabled() {
		logger.Info(
			"User exists with email",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.String("email", signupRequest.Email),
			zap.Bool("exists", userExists),
		)
	}

	if userExists {
		if logger.IsInfoEnabled() {
			logger.Info(
				"Email already registered",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.String("email", signupRequest.Email),
			)
		}
		return nil, utils.GetErrorResponse("Email already registered", 409)
	}

	hashedPassword, bcryptError := bcrypt.GenerateFromPassword([]byte(signupRequest.Password), 14)

	if bcryptError != nil {
		if logger.IsErrorEnabled() {
			logger.Error(
				"Error hashing password",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Error(bcryptError),
			)
		}
		return nil, utils.InternalServerErrorResponse()
	}

	signupRequest.Password = string(hashedPassword)
	dbUser := createUserEntity(signupRequest)
	dbUser.UserLoginType = constants.UserEntityLoginTypeEmailAndPassword
	user, saveError := authDao.SaveUser(requestId, dbUser)

	if saveError != nil {
		if logger.IsErrorEnabled() {
			logger.Error(
				"Error saving user in DB",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Int16(constants.ErrorCodeLogKey, saveError.ErrorCode),
				zap.Any(constants.ErrorMessageLogKey, saveError.Message),
			)
		}
		return nil, saveError
	}

	if user == nil {
		if logger.IsErrorEnabled() {
			logger.Error(
				"Failed to register user",
				zap.String(constants.RequestIdLogKey, requestId),
			)
		}
		return nil, utils.InternalServerErrorResponse()
	}

	if user.Email != nil {
		notificationService.SendSignupSuccessEmail(requestId, *user.Email, user.Name)
	}

	kafka_service.GetInstance().PushUserRegisteredEvent(requestId, user.Id)

	return &authModels.SignupResponse{
		Message:    "Signup successful! You can now explore all of the exciting and amazing features",
		StatusCode: 201,
	}, nil
}

// Logout Function to handle logout request and invalidates the jwt token
func Logout(requestId string, logoutRequest authModels.LogoutRequest) (*authModels.LogoutResponse, *authModels.ErrorResponse) {
	if logger.IsInfoEnabled() {
		logger.Info(
			"Processing Logout Request",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.Any("logoutRequest", logoutRequest),
		)
	}
	return &authModels.LogoutResponse{
		Message: "Logout successful",
	}, nil
}

// ValidateToken Function to handle validate token request and validates the jwt token
func ValidateToken(requestId string, validateTokenRequest authModels.ValidateTokenRequest) (*authModels.ValidateTokenResponse, *authModels.ErrorResponse) {
	if logger.IsDebugEnabled() {
		logger.Debug(
			"Processing Validate Token Request",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.Any("userId", validateTokenRequest.UserId),
			zap.Any("authToken", validateTokenRequest.AuthToken),
		)
	}

	tokenValidateResp, err := tokenService.GetInstance().ValidateJwtToken(requestId, validateTokenRequest.AuthToken, validateTokenRequest.UserId)

	if err != nil {
		if logger.IsErrorEnabled() {
			logger.Error(
				"Error validating JWT auth token",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Int16(constants.ErrorCodeLogKey, err.ErrorCode),
				zap.Any(constants.ErrorMessageLogKey, err.Message),
			)
		}
		return nil, err
	}

	return tokenValidateResp, nil
}

// GenerateAndSendForgotPasswordToken Function to generate forgot password token and send forgot password email back to user
func GenerateAndSendForgotPasswordToken(requestId string, forgotPasswordRequest authModels.ForgotPasswordRequest) (*authModels.ForgotPasswordResponse, *authModels.ErrorResponse) {
	if logger.IsDebugEnabled() {
		logger.Debug(
			"Processing forgot password Request",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.Any("forgotPasswordRequest", forgotPasswordRequest),
		)
	}

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

	if user.LoginType == constants.UserEntityLoginTypeOauthAndOtp || user.LoginType == constants.UserEntityLoginTypeOauthOnly {
		return nil, &authModels.ErrorResponse{
			Message:   "Invalid Request",
			Errors:    fmt.Sprintf("You are not allowed to reset password. Please login using %s oAuth", user.OAuthProvider),
			ErrorCode: 400,
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

// VerifyResetPasswordToken Function to validate forgot password token and return redirect URL to reset password UI page
func VerifyResetPasswordToken(requestId string, queryParams url.Values) (string, *authModels.ErrorResponse) {
	emailParam := queryParams["email"]
	tokenParam := queryParams["token"]

	if logger.IsInfoEnabled() {
		logger.Info("Processing reset password token request",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.Any("email", emailParam),
			zap.Any("token", tokenParam),
		)
	}

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
		if logger.IsErrorEnabled() {
			logger.Error(
				"Forgot Token not found in Database",
				zap.String(constants.RequestIdLogKey, requestId),
			)
		}

		return "", &authModels.ErrorResponse{
			Message:   "Invalid forgot password token. Please try again",
			ErrorCode: 400,
		}
	}

	redirectUrl := utils.GenerateForgotPasswordTokenRedirectUrl(email, token)

	if logger.IsDebugEnabled() {
		logger.Debug(
			"Redirect URL generated is",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.String("redirectUrl", redirectUrl),
		)
	}

	return redirectUrl, nil
}

// ResetPassword Function to actually reset password from forgot-password UI page
func ResetPassword(requestId string, resetPasswordRequest authModels.ResetPasswordRequest) (*authModels.ResetPasswordResponse, *authModels.ErrorResponse) {
	if logger.IsInfoEnabled() {
		logger.Info(
			"Processing Reset password Request",
			zap.String(constants.RequestIdLogKey, requestId),
		)
	}

	email := resetPasswordRequest.Email
	resetPasswordToken := resetPasswordRequest.ResetPasswordToken
	password := resetPasswordRequest.Password
	confirmPassword := resetPasswordRequest.ConfirmPassword

	// verify and match password
	if strings.TrimSpace(password) != strings.TrimSpace(confirmPassword) {
		if logger.IsErrorEnabled() {
			logger.Error(
				"Password & confirm Passwords mismatch",
				zap.String(constants.RequestIdLogKey, requestId),
			)
		}

		return nil, &authModels.ErrorResponse{
			Message:   "Password & confirm Passwords mismatch",
			ErrorCode: 400,
		}
	}

	// fetch forgot password token from DB
	forgotPasswordTokenFromDatabase, fptfdError := authDao.GetForgotPasswordToken(requestId, email)

	if fptfdError != nil {
		if logger.IsErrorEnabled() {
			logger.Error(
				"Failed to fetch forgot password token from DB",
				zap.String(constants.RequestIdLogKey, requestId),
			)
		}
		return nil, fptfdError
	}

	// match provided token with DB token again for double check
	if strings.TrimSpace(forgotPasswordTokenFromDatabase) != strings.TrimSpace(resetPasswordToken) {
		if logger.IsErrorEnabled() {
			logger.Error(
				"forgot token from DB doesn't match with token provided",
				zap.String(constants.RequestIdLogKey, requestId),
			)
		}

		return nil, &authModels.ErrorResponse{
			Message:   "Invalid token provided",
			ErrorCode: 400,
		}
	}

	hashedPassword, bcryptError := bcrypt.GenerateFromPassword([]byte(password), 14)

	if bcryptError != nil {
		if logger.IsErrorEnabled() {
			logger.Error(
				"Error hashing password",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Error(bcryptError),
			)
		}
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

// VerifyAdmin Function to check if userId is associated with an admin account or not
func VerifyAdmin(requestId string, verifyAdminRequest authModels.VerifyAdminRequest) (*authModels.VerifyAdminResponse, *authModels.ErrorResponse) {
	if logger.IsInfoEnabled() {
		logger.Info(
			"Processing Verify admin Request",
			zap.String(constants.RequestIdLogKey, requestId),
		)
	}

	userId := verifyAdminRequest.UserId

	user, err := authDao.GetUserById(requestId, userId)

	if err != nil {
		if logger.IsErrorEnabled() {
			logger.Error(
				"Failed to fetch admin user by ID",
				zap.String(constants.RequestIdLogKey, requestId),
			)
		}
		return nil, err
	}

	scopes := strings.Split(user.Scopes, ",")

	var adminScopeFound = false

	for _, scope := range scopes {
		if strings.Contains(strings.ToLower(scope), "admin") {
			adminScopeFound = true
			break
		}
	}

	if !adminScopeFound {
		response := &authModels.VerifyAdminResponse{
			Success:    false,
			Message:    "Admin scope not found",
			StatusCode: 200,
		}
		return response, nil
	}

	return &authModels.VerifyAdminResponse{
		Success:    true,
		Message:    "Admin verified successfully",
		StatusCode: 200,
	}, nil
}

// function to validate provided password against the encrypted password stored in DB
func verifyPassword(rawPassword string, encryptedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(rawPassword)) == nil
}

func createUserEntity(request model.SignupRequest) *entity.User {
	return &entity.User{
		Id:       strings.ReplaceAll(uuid.New().String(), "-", ""),
		Email:    &request.Email,
		Password: &request.Password,
		Name:     request.Name,
		Scopes:   "user",
	}
}
