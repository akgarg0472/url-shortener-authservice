package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	authModels "github.com/akgarg0472/urlshortener-auth-service/internal/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	utils "github.com/akgarg0472/urlshortener-auth-service/utils"
)

var logger = Logger.GetLogger("requestBodyValidator.go")

func LoginRequestBodyValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Header.Get("Request-ID")

		var loginRequest authModels.LoginRequest

		decodeError := decodeRequestBody(r, &loginRequest)

		if decodeError != nil {
			logger.Error("[{}]: Error decoding request body: {}", requestId, decodeError.Error())
			resp := generateErrorResponse("Invalid request body", 400)
			writeErrorResponse(w, http.StatusBadRequest, resp)
			return
		}

		validationErrors := utils.ValidateRequestFields(loginRequest)

		if validationErrors != nil {
			logger.Error("[{}]: Login Request Validation failed: {}", requestId, validationErrors)
			errorResponse, _ := json.Marshal(validationErrors)
			writeErrorResponse(w, http.StatusBadRequest, errorResponse)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "loginRequest", loginRequest)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func SignupRequestBodyValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Header.Get("Request-ID")

		var signupRequest authModels.SignupRequest

		decodeError := decodeRequestBody(r, &signupRequest)

		if decodeError != nil {
			logger.Error("[{}]: Error decoding request body: {}", requestId, decodeError.Error())
			resp := generateErrorResponse("Invalid request body", 400)
			writeErrorResponse(w, http.StatusBadRequest, resp)
			return
		}

		validationErrors := utils.ValidateRequestFields(signupRequest)

		if validationErrors != nil {
			logger.Error("[{}]: Signup Request Validation failed")
			errorResponse, _ := json.Marshal(validationErrors)
			writeErrorResponse(w, http.StatusBadRequest, errorResponse)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "signupRequest", signupRequest)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func decodeRequestBody(r *http.Request, ref interface{}) error {
	return json.NewDecoder(r.Body).Decode(&ref)
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(message)
}

func generateErrorResponse(message string, errorCode int16) []byte {
	errorResponse := authModels.ErrorResponse{
		Message:   message,
		ErrorCode: int(errorCode),
	}
	resp, _ := json.Marshal(errorResponse)
	return resp
}
