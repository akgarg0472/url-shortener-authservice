package handler

import (
	"net/http"

	AuthModels "github.com/akgarg0472/urlshortener-auth-service/model"
	utils "github.com/akgarg0472/urlshortener-auth-service/utils"
)

// Function to send response back to client
func sendResponseToClient(responseWriter http.ResponseWriter, requestId string, response interface{}, err *AuthModels.ErrorResponse, statusCode int) {
	if err != nil {
		errorJson, _ := utils.ConvertToJsonString(err)
		sendResponseToClientWithStatusAndMessage(responseWriter, int(err.ErrorCode), errorJson)
		return
	}

	jsonResponse, jsonConvertError := utils.ConvertToJsonString(response)

	if jsonConvertError != nil {
		authLogger.Error("[{}]: Error Converting Response to JSON: {}", requestId, jsonConvertError.Error())

		errorResponse := &AuthModels.ErrorResponse{
			Message:   "Internal Server Error",
			ErrorCode: 500,
		}

		errorResponseJson, _ := utils.ConvertToJsonString(errorResponse)
		sendResponseToClientWithStatusAndMessage(responseWriter, http.StatusInternalServerError, errorResponseJson)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)
	_, err2 := responseWriter.Write([]byte(jsonResponse))
	if err2 != nil {
		return
	}
}

// Function to send response to client with given status code and message
func sendResponseToClientWithStatusAndMessage(responseWriter http.ResponseWriter, statusCode int, message string) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)
	_, err := responseWriter.Write([]byte(message))
	if err != nil {
		return
	}
}
