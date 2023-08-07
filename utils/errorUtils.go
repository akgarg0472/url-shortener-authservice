package utils

import (
	"fmt"
	"regexp"
	"strings"

	AuthModels "github.com/akgarg0472/urlshortener-auth-service/model"
	"github.com/go-sql-driver/mysql"
)

func InternalServerErrorResponseByte() []byte {
	errorResponse := AuthModels.ErrorResponse{
		Message:   "internal server error",
		ErrorCode: 500,
	}

	errorResponseJson, _ := ConvertToJsonString(errorResponse)
	return []byte(errorResponseJson)
}

func InternalServerErrorResponse() *AuthModels.ErrorResponse {
	return &AuthModels.ErrorResponse{
		Message:   "internal server error",
		ErrorCode: 500,
	}
}

func GetErrorResponse(message interface{}, errorCode int16) *AuthModels.ErrorResponse {
	return &AuthModels.ErrorResponse{
		Message:   message,
		ErrorCode: errorCode,
	}
}

func ParseMySQLErrorAndReturnErrorResponse(err error) *AuthModels.ErrorResponse {
	mysqlErr, ok := err.(*mysql.MySQLError)

	if !ok {
		return InternalServerErrorResponse()
	}

	switch mysqlErr.Number {
	case 1062:
		key, value := extractDuplicatedKeyAndValue(mysqlErr.Message)
		return GetErrorResponse(fmt.Sprintf("%s '%s' already exists", key, value), 409)
	}

	return InternalServerErrorResponse()
}

func extractDuplicatedKeyAndValue(errorMessage string) (string, string) {
	re := regexp.MustCompile(`Duplicate entry '(.*)' for key '(.*)'`)
	matches := re.FindStringSubmatch(errorMessage)

	if len(matches) == 3 {
		key := strings.Split(matches[2], ".")
		return key[len(key)-1], matches[1]
	}

	return "null", "null"
}
