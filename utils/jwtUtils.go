package utils

import (
	"time"

	Model "github.com/akgarg0472/urlshortener-auth-service/model"
	"github.com/golang-jwt/jwt/v5"
)

func ParseAndGenerateJwtErrorResponse(claims jwt.MapClaims) *Model.ErrorResponse {
	expirationTime := claims["exp"]

	if isTokenExpired(int64(expirationTime.(float64))) {
		return &Model.ErrorResponse{
			Message:   "JWT_TOKEN_EXPIRED",
			ErrorCode: 400,
		}
	}

	return &Model.ErrorResponse{
		Message:   "JWT_TOKEN_INVALID",
		ErrorCode: 400,
	}
}

func isTokenExpired(expirationTimestamp int64) bool {
	return time.Now().After(time.Unix(expirationTimestamp, 0))
}
