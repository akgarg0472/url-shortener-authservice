package oauth_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/akgarg0472/urlshortener-auth-service/constants"
	"github.com/akgarg0472/urlshortener-auth-service/internal/logger"
	"github.com/akgarg0472/urlshortener-auth-service/model"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
	"go.uber.org/zap"
)

const (
	GoogleAccessTokenUrl = "https://oauth2.googleapis.com/token"
	GoogleUserInfoUrl    = "https://www.googleapis.com/oauth2/v2/userinfo"
)

type GoogleAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type GoogleUserInfoResponse struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	ProfilePicture string `json:"picture"`
	Email          string `json:"email"`
}

func FetchGoogleProfileInfo(requestId string, request model.OAuthCallbackRequest) (*ProfileInfo, *model.ErrorResponse) {
	if logger.IsInfoEnabled() {
		logger.Info(
			"Fetching profile info from google",
			zap.String(constants.RequestIdLogKey, requestId),
		)
	}

	clientId := utils.GetEnvVariable("OAUTH_GOOGLE_CLIENT_ID", "")
	clientSecret := utils.GetEnvVariable("OAUTH_GOOGLE_CLIENT_SECRET", "")
	redirectUri := utils.GetEnvVariable("OAUTH_GOOGLE_CLIENT_REDIRECT_URI", "")

	requestBody, err := json.Marshal(map[string]string{
		"code":          request.Code,
		"client_id":     clientId,
		"client_secret": clientSecret,
		"grant_type":    "authorization_code",
		"redirect_uri":  redirectUri,
	})

	if err != nil {
		if logger.IsErrorEnabled() {
			logger.Error(
				"Error creating JSON request",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Error(err),
			)
		}
		return nil, utils.InternalServerErrorResponse()
	}

	resp, err := http.Post(GoogleAccessTokenUrl, "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		if logger.IsErrorEnabled() {
			logger.Error(
				"Failed to get access token",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Error(err),
			)
		}
		return nil, utils.InternalServerErrorResponse()
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// do nothing
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var response map[string]any
		er := json.NewDecoder(resp.Body).Decode(&response)

		if er == nil {
			if logger.IsErrorEnabled() {
				logger.Error(
					"Non 2xx status code received from access token request",
					zap.String(constants.RequestIdLogKey, requestId),
					zap.Int(constants.StatusCodeLogKey, resp.StatusCode),
					zap.Any("response_body", response),
				)
			}
		} else {
			if logger.IsErrorEnabled() {
				logger.Error(
					"Non 2xx status code received from access token request",
					zap.String(constants.RequestIdLogKey, requestId),
					zap.Int(constants.StatusCodeLogKey, resp.StatusCode),
				)
			}
		}
		return nil, utils.InternalServerErrorResponse()
	}

	var tokenResponse GoogleAccessTokenResponse

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		if logger.IsErrorEnabled() {
			logger.Error(
				"Failed to decode access token response",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Error(err),
			)
		}
		return nil, utils.InternalServerErrorResponse()
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", GoogleUserInfoUrl, nil)

	if err != nil {
		if logger.IsErrorEnabled() {
			logger.Error(
				"Failed to create user info request",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Error(err),
			)
		}
		return nil, utils.InternalServerErrorResponse()
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", tokenResponse.TokenType, tokenResponse.AccessToken))

	resp, err = client.Do(req)

	if err != nil {
		if logger.IsErrorEnabled() {
			logger.Error(
				"Failed to get user info",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Error(err),
			)
		}
		return nil, utils.InternalServerErrorResponse()
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// do nothing
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		if logger.IsErrorEnabled() {
			logger.Error(
				"Invalid status code received from user info request",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Int(constants.StatusCodeLogKey, resp.StatusCode),
			)
		}
		return nil, utils.InternalServerErrorResponse()
	}

	var userInfoResponse GoogleUserInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&userInfoResponse); err != nil {
		if logger.IsErrorEnabled() {
			logger.Error(
				"Failed to decode user info response",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Error(err),
			)
		}
		return nil, utils.InternalServerErrorResponse()
	}

	return &ProfileInfo{
		OAuthId:        userInfoResponse.Id,
		Name:           userInfoResponse.Name,
		ProfilePicture: userInfoResponse.ProfilePicture,
		Email:          userInfoResponse.Email,
	}, nil
}
