package oauth_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/akgarg0472/urlshortener-auth-service/constants"
	"github.com/akgarg0472/urlshortener-auth-service/internal/logger"
	"github.com/akgarg0472/urlshortener-auth-service/model"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
	"go.uber.org/zap"
)

const (
	GithubAccessTokenUrl = "https://github.com/login/oauth/access_token"
	GithubUserInfoUrl    = "https://api.github.com/user"
)

type GitHubAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type GitHubUserInfoResponse struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	ProfilePicture string `json:"avatar_url"`
	Email          string `json:"email"`
}

func FetchGitHubProfileInfo(requestId string, request model.OAuthCallbackRequest) (*ProfileInfo, *model.ErrorResponse) {
	if logger.IsInfoEnabled() {
		logger.Info(
			"fetching profile info from GitHub",
			zap.String(constants.RequestIdLogKey, requestId),
		)
	}

	clientId := utils.GetEnvVariable("OAUTH_GITHUB_CLIENT_ID", "")
	clientSecret := utils.GetEnvVariable("OAUTH_GITHUB_CLIENT_SECRET", "")
	redirectUri := utils.GetEnvVariable("OAUTH_GITHUB_CLIENT_REDIRECT_URI", "")

	requestBody, err := json.Marshal(map[string]string{
		"code":          request.Code,
		"client_id":     clientId,
		"client_secret": clientSecret,
		"redirect_uri":  redirectUri,
	})

	if err != nil {
		if logger.IsErrorEnabled() {
			logger.Error(
				"error creating JSON request",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Error(err),
			)
		}
		return nil, utils.InternalServerErrorResponse()
	}

	client := &http.Client{}

	req, _ := http.NewRequest("POST", GithubAccessTokenUrl, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		if logger.IsErrorEnabled() {
			logger.Error(
				"failed to get access token from GitHub",
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
				"invalid status code received from access token request",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Int(constants.StatusCodeLogKey, resp.StatusCode),
			)
		}
		return nil, utils.InternalServerErrorResponse()
	}

	if logger.IsInfoEnabled() {
		logger.Info(
			"Access token successfully fetched",
			zap.String(constants.RequestIdLogKey, requestId),
		)
	}

	var tokenResponse GitHubAccessTokenResponse

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		if logger.IsErrorEnabled() {
			logger.Error(
				"Error decoding access token response",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Error(err),
			)
		}
		return nil, utils.InternalServerErrorResponse()
	}

	req, err = http.NewRequest("GET", GithubUserInfoUrl, nil)

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
				"Failed to retrieve user info from GitHub",
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

	var userInfoResponse GitHubUserInfoResponse
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
		OAuthId:        strconv.FormatInt(userInfoResponse.Id, 10),
		Name:           userInfoResponse.Name,
		ProfilePicture: userInfoResponse.ProfilePicture,
		Email:          userInfoResponse.Email,
	}, nil
}
