package token_service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/akgarg0472/urlshortener-auth-service/constants"
	"github.com/akgarg0472/urlshortener-auth-service/internal/logger"
	"github.com/akgarg0472/urlshortener-auth-service/model"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

var (
	instance *TokenService
)

type TokenService struct {
	jwtSecretKey            []byte
	jwtIssuer               string
	jwtValidity             int64
	forgotPasswordSecretKey []byte
	forgotPasswordValidity  int64
}

func GetInstance() *TokenService {
	if instance == nil {
		instance = &TokenService{
			jwtSecretKey:            []byte(getJWTSecretKey()),
			jwtIssuer:               getJWTIssuer(),
			jwtValidity:             getJWTValidityDurationInSeconds(),
			forgotPasswordSecretKey: []byte(getForgotPasswordSecretKey()),
			forgotPasswordValidity:  getForgotPasswordValidityDurationInSeconds(),
		}
	}

	return instance
}

// GenerateJwtToken generates the JWT token and stores it in the map
func (tokenService *TokenService) GenerateJwtToken(requestId string, user model.User) (string, *model.ErrorResponse) {
	if logger.IsInfoEnabled() {
		logger.Info("Generating JWT token",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.String("user", user.String()),
		)
	}

	claims := jwt.MapClaims{
		"iss":    tokenService.jwtIssuer,
		"sub":    user.Email,
		"uid":    user.Id,
		"scopes": user.Scopes,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Unix() + tokenService.jwtValidity,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtTokenString, err := token.SignedString(tokenService.jwtSecretKey)

	if err != nil {
		if logger.IsErrorEnabled() {
			logger.Error("Error while generating JWT token",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Error(err),
			)
		}
		return "", utils.InternalServerErrorResponse()
	}

	return jwtTokenString, nil
}

// ValidateJwtToken validates the JWT token by checking if it is valid and not expired
func (tokenService *TokenService) ValidateJwtToken(
	requestId string,
	jwtToken string,
	userId string,
) (*model.ValidateTokenResponse, *model.ErrorResponse) {
	if logger.IsDebugEnabled() {
		logger.Debug("Validating JWT token",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.String("userId", userId),
			zap.String("jwtToken", jwtToken),
		)
	}

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error parsing token")
		}
		return tokenService.jwtSecretKey, nil
	})

	if err != nil {
		if logger.IsErrorEnabled() {
			logger.Error("Error validating token",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Error(err),
			)
		}

		if token != nil {
			claims, isMapClaims := token.Claims.(jwt.MapClaims)

			if isMapClaims {
				return nil, utils.ParseAndGenerateJwtErrorResponse(claims)
			}
		}

		return nil, utils.BadRequestErrorResponse("JWT_TOKEN_INVALID")
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	uId := claims["uid"].(string)

	if strings.TrimSpace(uId) != strings.TrimSpace(userId) {
		if logger.IsErrorEnabled() {
			logger.Error("Error validating token: Invalid userId",
				zap.String(constants.RequestIdLogKey, requestId),
			)
		}
		return nil, utils.BadRequestErrorResponse("Passwords mismatch")
	}

	return &model.ValidateTokenResponse{
		UserId:     uId,
		Expiration: claims["exp"].(float64),
		Token:      token.Raw,
		Success:    userId == uId,
	}, nil
}

func (tokenService *TokenService) GenerateForgotPasswordToken(requestId string, user model.User) (string, *model.ErrorResponse) {
	if logger.IsInfoEnabled() {
		logger.Info("Generating forgot password token",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.String("user", user.String()),
		)
	}

	claims := jwt.MapClaims{
		"sub": user.Email,
		"iat": time.Now().Unix(),
		"exp": time.Now().Unix() + tokenService.forgotPasswordValidity,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	forgotPasswordToken, err := token.SignedString(tokenService.forgotPasswordSecretKey)

	if err != nil {
		if logger.IsErrorEnabled() {
			logger.Error("Error while generating forgot password token",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Error(err),
			)
		}
		return "", utils.InternalServerErrorResponse()
	}

	return forgotPasswordToken, nil
}

func (tokenService *TokenService) ValidateForgotPasswordToken(requestId string, forgotPasswordToken string) *model.ErrorResponse {
	if logger.IsDebugEnabled() {
		logger.Debug("Validating Forgot Password token",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.String("forgotPasswordToken", forgotPasswordToken),
		)
	}

	parsedToken, err := jwt.Parse(forgotPasswordToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error parsing token")
		}
		return tokenService.forgotPasswordSecretKey, nil
	})

	if err != nil {
		if logger.IsErrorEnabled() {
			logger.Error("Error validating token",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Error(err),
			)
		}

		if parsedToken != nil {
			claims, isMapClaims := parsedToken.Claims.(jwt.MapClaims)

			if isMapClaims {
				return utils.ParseAndGenerateJwtErrorResponse(claims)
			}
		}

		return utils.BadRequestErrorResponse("JWT_TOKEN_INVALID")
	}

	return nil
}

func getJWTSecretKey() string {
	secret := utils.GetEnvVariable("JWT_SECRET_KEY", "")

	if secret == "" {
		panic("JWT_SECRET_KEY not found")
	}

	return secret
}

func getJWTIssuer() string {
	return utils.GetEnvVariable("JWT_TOKEN_ISSUER", "auth-service")
}

func getJWTValidityDurationInSeconds() int64 {
	expiry := utils.GetEnvVariable("JWT_TOKEN_EXPIRY", "3600000")

	value, err := strconv.ParseInt(expiry, 10, 64)

	if err != nil {
		return 3600000
	} else {
		return value
	}
}

func getForgotPasswordSecretKey() string {
	secret := utils.GetEnvVariable("FORGOT_PASS_SECRET_KEY", "")

	if secret == "" {
		panic("FORGOT_PASS_SECRET_KEY not found")
	}

	return secret
}

func getForgotPasswordValidityDurationInSeconds() int64 {
	expiry := utils.GetEnvVariable("FORGOT_PASS_EXPIRY", "600")

	value, err := strconv.ParseInt(expiry, 10, 64)

	if err != nil {
		return 600
	} else {
		return value
	}
}
