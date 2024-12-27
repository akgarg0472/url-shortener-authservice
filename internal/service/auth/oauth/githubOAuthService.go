package oauth_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
)

var (
	gHOALogger = Logger.GetLogger("githubOAuthService.go")
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

func FetchGitHubProfileInfo(reqId string, request model.OAuthCallbackRequest) (*ProfileInfo, *model.ErrorResponse) {
	gHOALogger.Info("[{}] fetching profile info from GitHub", reqId)

	clientId := utils.GetEnvVariable("OAUTH_GITHUB_CLIENT_ID	", "")
	clientSecret := utils.GetEnvVariable("OAUTH_GITHUB_CLIENT_SECRET", "")
	redirectUri := utils.GetEnvVariable("OAUTH_GITHUB_CLIENT_REDIRECT_URI", "")

	requestBody, err := json.Marshal(map[string]string{
		"code":          request.Code,
		"client_id":     clientId,
		"client_secret": clientSecret,
		"redirect_uri":  redirectUri,
	})

	if err != nil {
		gHOALogger.Error("[{}] error creating JSON request: {}", reqId, err)
		return nil, utils.InternalServerErrorResponse()
	}

	client := &http.Client{}

	req, _ := http.NewRequest("POST", GithubAccessTokenUrl, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		gHOALogger.Error("[{}] failed to get access token: {}", reqId, err)
		return nil, utils.InternalServerErrorResponse()
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// do nothing
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		gHOALogger.Error("[{}] invalid status code received from access token request: {}", reqId, resp.StatusCode)
		return nil, utils.InternalServerErrorResponse()
	}

	gHOALogger.Info("[{}] access token successfully fetched", reqId)

	var tokenResponse GitHubAccessTokenResponse

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		gHOALogger.Error("[{}] failed to decode token response: {}", reqId, err)
		return nil, utils.InternalServerErrorResponse()
	}

	req, err = http.NewRequest("GET", GithubUserInfoUrl, nil)

	if err != nil {
		gHOALogger.Error("[{}] failed to create user info request: {}", reqId, err)
		return nil, utils.InternalServerErrorResponse()
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", tokenResponse.TokenType, tokenResponse.AccessToken))

	resp, err = client.Do(req)

	if err != nil {
		gHOALogger.Error("[{}] failed to get user info: {}", reqId, err)
		return nil, utils.InternalServerErrorResponse()
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// do nothing
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		gHOALogger.Error("[{}] invalid status code received from user info request: {}", reqId, resp.StatusCode)
		return nil, utils.InternalServerErrorResponse()
	}

	var userInfoResponse GitHubUserInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&userInfoResponse); err != nil {
		gHOALogger.Error("[{}] failed to decode user info response: {}", reqId, err)
		return nil, utils.InternalServerErrorResponse()
	}

	return &ProfileInfo{
		OAuthId:        strconv.FormatInt(userInfoResponse.Id, 10),
		Name:           userInfoResponse.Name,
		ProfilePicture: userInfoResponse.ProfilePicture,
		Email:          userInfoResponse.Email,
	}, nil
}
