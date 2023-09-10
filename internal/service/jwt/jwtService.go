package jwt_service

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
	logger   = Logger.GetLogger("jwtService.go")
	instance *JwtService
)

type JwtService struct {
	secretKey []byte
	issuer    string
	validity  int64
}

func GetInstance() *JwtService {
	if instance == nil {
		instance = &JwtService{
			secretKey: []byte(getSecretKey()),
			issuer:    getJwtIssuer(),
			validity:  getJwtValidityDurationInSeconds(),
		}
	}

	return instance
}

// GenerateJwtToken generates the JWT token and stores it in the map
func (jwtService *JwtService) GenerateJwtToken(requestId string, user Model.User) (string, *Model.ErrorResponse) {
	logger.Info("[{}]: Generating JWT token -> {}", requestId, user.String())

	claims := jwt.MapClaims{
		"iss":    jwtService.issuer,
		"sub":    user.Email,
		"uid":    user.Id,
		"scopes": user.Scopes,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Unix() + jwtService.validity,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtTokenString, err := token.SignedString([]byte(jwtService.secretKey))

	if err != nil {
		logger.Error("[{}]: Error while generating JWT token -> {}", requestId, err.Error())
		return "", utils.InternalServerErrorResponse()
	}

	return jwtTokenString, nil
}

// ValidateJwtToken validates the JWT token by checking if it exists in the map and is not expired
func (jwtService *JwtService) ValidateJwtToken(requestId string, jwtToken string, userId string) (*Model.ValidateTokenResponse, *Model.ErrorResponse) {
	logger.Debug("[{}]: Validating JWT token -> {}, {}", requestId, userId, jwtToken)

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error parsing token")
		}
		return jwtService.secretKey, nil
	})

	if err != nil {
		logger.Error("[{}]: Error validating token -> {}", requestId, err.Error())

		claims, isMapClaims := token.Claims.(jwt.MapClaims)

		if isMapClaims {
			return nil, utils.ParseAndGenerateJwtErrorResponse(claims)
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

func getSecretKey() string {
	secret := utils.GetEnvVariable("JWT_SECRET_KEY", "")

	if secret == "" {
		panic("JWT_SECRET_KEY not found")
	}

	return secret
}

func getJwtIssuer() string {
	return utils.GetEnvVariable("JWT_TOKEN_ISSUER", "auth-service")
}

func getJwtValidityDurationInSeconds() int64 {
	expiry := utils.GetEnvVariable("JWT_TOKEN_EXPIRY", "3600000")

	value, err := strconv.ParseInt(expiry, 10, 64)

	if err != nil {
		return 3600000
	} else {
		return value
	}
}
