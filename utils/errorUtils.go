package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/akgarg0472/urlshortener-auth-service/model"
	"github.com/go-sql-driver/mysql"
)

func InternalServerErrorResponseByte() []byte {
	errorResponse := model.ErrorResponse{
		Message:   "internal server error",
		ErrorCode: 500,
	}

	errorResponseJson, _ := ConvertToJsonString(errorResponse)
	return []byte(errorResponseJson)
}

func InternalServerErrorResponse() *model.ErrorResponse {
	return &model.ErrorResponse{
		Message:   "internal server error",
		ErrorCode: 500,
	}
}

func BadRequestErrorResponse(message string) *model.ErrorResponse {
	return &model.ErrorResponse{
		Message:   message,
		ErrorCode: 400,
	}
}

func GetErrorResponse(message interface{}, errorCode int16) *model.ErrorResponse {
	return &model.ErrorResponse{
		Message:   message,
		ErrorCode: errorCode,
	}
}

func GetErrorResponseByte(message interface{}, errorCode int16) []byte {
	resp := GetErrorResponse(message, errorCode)
	errorResponseJson, _ := ConvertToJsonBytes(resp)
	return errorResponseJson
}

func ParseMySQLErrorAndReturnErrorResponse(err error) *model.ErrorResponse {
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
