package utils

import (
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
