package token_service

import (
	"fmt"
	"strconv"
	"time"

	Model "github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
	"github.com/golang-jwt/jwt/v5"
)

var (
	logger   = Logger.GetLogger("tokenService.go")
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
func (tokenService *TokenService) GenerateJwtToken(requestId string, user Model.User) (string, *Model.ErrorResponse) {
	logger.Info("[{}]: Generating JWT token -> {}", requestId, user.String())

	claims := jwt.MapClaims{
		"iss":    tokenService.jwtIssuer,
		"sub":    user.Email,
		"uid":    user.Id,
		"scopes": user.Scopes,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Unix() + tokenService.jwtValidity,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtTokenString, err := token.SignedString([]byte(tokenService.jwtSecretKey))

	if err != nil {
		logger.Error("[{}]: Error while generating JWT token -> {}", requestId, err.Error())
		return "", utils.InternalServerErrorResponse()
	}

	return jwtTokenString, nil
}

// ValidateJwtToken validates the JWT token by checking if it exists in the map and is not expired
func (tokenService *TokenService) ValidateJwtToken(requestId string, jwtToken string, userId string) (*Model.ValidateTokenResponse, *Model.ErrorResponse) {
	logger.Debug("[{}]: Validating JWT token -> {}, {}", requestId, userId, jwtToken)

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error parsing token")
		}
		return tokenService.jwtSecretKey, nil
	})

	if err != nil {
		logger.Error("[{}]: Error validating token -> {}", requestId, err.Error())

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

	return &Model.ValidateTokenResponse{
		UserId:     uId,
		Expiration: claims["exp"].(float64),
		Token:      token.Raw,
		Success:    userId == uId,
	}, nil
}

func (tokenService *TokenService) GenerateForgotPasswordToken(requestId string, user Model.User) (string, *Model.ErrorResponse) {
	logger.Info("[{}]: Generating forgot password token -> {}", requestId, user.String())

	claims := jwt.MapClaims{
		"sub": user.Email,
		"iat": time.Now().Unix(),
		"exp": time.Now().Unix() + tokenService.forgotPasswordValidity,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	forgotPasswordToken, err := token.SignedString([]byte(tokenService.forgotPasswordSecretKey))

	if err != nil {
		logger.Error("[{}]: Error while generating forgot password token -> {}", requestId, err.Error())
		return "", utils.InternalServerErrorResponse()
	}

	return forgotPasswordToken, nil
}

func (tokenService *TokenService) ValidateForgotPasswordToken(requestId string, forgotPasswordToken string) *Model.ErrorResponse {
	logger.Debug("[{}]: Validating Forgot Password token -> {}", requestId, forgotPasswordToken)

	parsedToken, err := jwt.Parse(forgotPasswordToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error parsing token")
		}
		return tokenService.forgotPasswordSecretKey, nil
	})

	if err != nil {
		logger.Error("[{}]: Error validating token -> {}", requestId, err.Error())

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
