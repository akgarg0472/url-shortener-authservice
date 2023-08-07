package jwt_service

import (
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
	secretKey      []byte
	issuer         string
	validityMillis int64
}

func GetInstance() *JwtService {
	if instance == nil {
		instance = &JwtService{
			secretKey:      []byte(getSecretKey()),
			issuer:         getJwtIssuer(),
			validityMillis: getJwtValidityDurationInMillis(),
		}
	}

	return instance
}

func (jwtService *JwtService) GenerateJwtToken(requestId string, user Model.User) (string, *Model.ErrorResponse) {
	logger.Trace("[{}]: Generating JWT token -> {}", requestId, user.String())

	claims := jwt.MapClaims{
		"iss":    jwtService.issuer,
		"sub":    user.Email,
		"userid": user.Id,
		"scopes": user.Scopes,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().UnixMilli() + jwtService.validityMillis,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtTokenString, err := token.SignedString([]byte(jwtService.secretKey))

	if err != nil {
		logger.Error("[{}]: Error while generating JWT token -> {}", requestId, err.Error())
		return "", utils.InternalServerErrorResponse()
	}

	return jwtTokenString, nil
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

func getJwtValidityDurationInMillis() int64 {
	expiry := utils.GetEnvVariable("JWT_TOKEN_EXPIRY", "3600000")

	value, err := strconv.ParseInt(expiry, 10, 64)

	if err != nil {
		return 3600000
	} else {
		return value
	}
}
