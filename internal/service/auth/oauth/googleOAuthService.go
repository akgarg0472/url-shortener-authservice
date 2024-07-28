package oauth_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
)

var (
	gOALogger = Logger.GetLogger("googleOAuthService.go")
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

func FetchGoogleProfileInfo(reqId string, request model.OAuthCallbackRequest) (*ProfileInfo, *model.ErrorResponse) {
	gOALogger.Info("[{}] fetching profile info from google", reqId)

	clientId := utils.GetEnvVariable("OAUTH_CLIENT_GOOGLE_ID", "")
	clientSecret := utils.GetEnvVariable("OAUTH_CLIENT_GOOGLE_SECRET", "")
	redirectUri := utils.GetEnvVariable("OAUTH_CLIENT_GOOGLE_REDIRECT_URI", "")

	requestBody, err := json.Marshal(map[string]string{
		"code":          request.Code,
		"client_id":     clientId,
		"client_secret": clientSecret,
		"grant_type":    "authorization_code",
		"redirect_uri":  redirectUri,
	})

	if err != nil {
		logger.Error("[{}] error creating JSON request: {}", reqId, err)
		return nil, utils.InternalServerErrorResponse()
	}

	resp, err := http.Post(GoogleAccessTokenUrl, "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		logger.Error("[{}] failed to get access token: {}", reqId, err)
		return nil, utils.InternalServerErrorResponse()
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// do nothing
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		logger.Error("[{}] invalid status code received from access token request: {}", reqId, resp.StatusCode)
		return nil, utils.InternalServerErrorResponse()
	}

	var tokenResponse GoogleAccessTokenResponse

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		logger.Error("[{}] failed to decode token response: {}", reqId, err)
		return nil, utils.InternalServerErrorResponse()
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", GoogleUserInfoUrl, nil)

	if err != nil {
		logger.Error("[{}] failed to create user info request: {}", reqId, err)
		return nil, utils.InternalServerErrorResponse()
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", tokenResponse.TokenType, tokenResponse.AccessToken))

	resp, err = client.Do(req)

	if err != nil {
		logger.Error("[{}] failed to get user info: {}", reqId, err)
		return nil, utils.InternalServerErrorResponse()
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// do nothing
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		logger.Error("[{}] invalid status code received from user info request: {}", reqId, resp.StatusCode)
		return nil, utils.InternalServerErrorResponse()
	}

	var userInfoResponse GoogleUserInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&userInfoResponse); err != nil {
		logger.Error("[{}] failed to decode user info response: {}", reqId, err)
		return nil, utils.InternalServerErrorResponse()
	}

	return &ProfileInfo{
		OAuthId:        userInfoResponse.Id,
		Name:           userInfoResponse.Name,
		ProfilePicture: userInfoResponse.ProfilePicture,
		Email:          userInfoResponse.Email,
	}, nil
}
