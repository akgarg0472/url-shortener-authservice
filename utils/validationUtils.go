package utils

import (
	AuthModels "github.com/akgarg0472/urlshortener-auth-service/model"
	Validator "github.com/go-playground/validator/v10"
)

var validator = Validator.New()

func ValidateRequestFields(request interface{}) map[string]string {
	err := validator.Struct(request)

	if err != nil {
		errorsMap := make(map[string]string)

		for _, error := range err.(Validator.ValidationErrors) {
			field := error.Field()
			tag := error.ActualTag()
			errorsMap[field] = convertValidationErrorToMessage(field, tag)
		}

		return errorsMap
	}

	return nil
}

func convertValidationErrorToMessage(field string, tag string) string {
	switch tag {
	case "required":
		return field + " is required"

	default:
		return field + " is invalid"
	}
}

func ValidateResetPasswordRequestQueryParams(emailParam []string, tokenParam []string) *AuthModels.ErrorResponse {
	emailParamLength := len(emailParam)
	tokenParamLength := len(tokenParam)

	if emailParamLength != 1 || tokenParamLength != 1 {
		var requestErrors = ""

		if emailParamLength != 1 {
			requestErrors += "'email' is missing, "
		}

		if tokenParamLength != 1 {
			requestErrors += "'token' is missing"
		}

		return &AuthModels.ErrorResponse{
			Message:   "Invalid Reset Password Request. Missing required query parameters",
			ErrorCode: 400,
			Errors:    requestErrors,
		}
	}

	return nil
}
