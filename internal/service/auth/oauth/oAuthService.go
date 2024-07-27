package oauth_service

import (
	"fmt"
	"strings"

	authDao "github.com/akgarg0472/urlshortener-auth-service/internal/dao/auth"
	"github.com/akgarg0472/urlshortener-auth-service/internal/dao/entity"
	oauthDao "github.com/akgarg0472/urlshortener-auth-service/internal/dao/oauth"
	notificationService "github.com/akgarg0472/urlshortener-auth-service/internal/service/notification"
	tokenService "github.com/akgarg0472/urlshortener-auth-service/internal/service/token"
	"github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
	"github.com/google/uuid"
)

var (
	oAuthClientMapping map[string]model.OAuthClient
	logger             = Logger.GetLogger("oAuthService.go")
)

type ProfileInfo struct {
	Id             string
	Name           string
	ProfilePicture string
	Email          string
	Username       string
}

func (p ProfileInfo) String() string {
	return fmt.Sprintf(`{"Id":"%s", "Name":"%s", "ProfilePicture":"%s", "Email":"%s", "Username":"%s"}`, p.Id, p.Name, p.ProfilePicture, p.Email, p.Username)
}

type AccessTokenResponse struct {
	AccessToken string
	TokenType   string
}

func InitOAuthClients() {
	clients := oauthDao.FetchOAuthClients()

	if oAuthClientMapping == nil {
		oAuthClientMapping = make(map[string]model.OAuthClient)
	}

	for _, _client := range clients {
		client := model.OAuthClient{
			Provider:    _client.Provider,
			ClientId:    _client.ClientID,
			RedirectURI: _client.RedirectURI,
			AccessType:  _client.AccessType,
			Scope:       _client.Scope,
		}

		oAuthClientMapping[client.Provider] = client
	}

	logger.Debug("Loaded oAuth clients: {}", oAuthClientMapping)
}

func GetOAuthClient(query string) []model.OAuthClient {
	providers := strings.Split(query, ",")

	clients := []model.OAuthClient{}

	for _, provider := range providers {
		if client, exists := oAuthClientMapping[provider]; exists {
			clients = append(clients, client)
		}
	}

	return clients
}

func ProcessCallbackRequest(requestId string, oAuthCallbackRequest model.OAuthCallbackRequest) (*model.OAuthCallbackResponse, *model.ErrorResponse) {
	logger.Info("[{}] oAuth callback request received: {}", requestId, oAuthCallbackRequest)

	profileInfo, err := getProfileInfo(requestId, oAuthCallbackRequest)

	if err != nil {
		logger.Error("[{}] error fetching oAuth profile details: {}", requestId, err)
		return nil, err
	}

	// checks if user is registered or not
	identifier := getIdentifier(*profileInfo)
	user, err := authDao.GetUserByEmailOrId(requestId, identifier)

	if user != nil {
		logger.Info("[{}] user is already registered", requestId)
	} else if user == nil && err.ErrorCode == 404 {
		logger.Info("[{}] user is not registered, going to register it", requestId)
		registeredUser, err := registerUser(requestId, *profileInfo)

		if err != nil {
			logger.Error("[{}] failed to register OAuth user: {}", requestId, err)
			return nil, err
		}

		user = registeredUser

		if user.Email != "" {
			notificationService.SendSignupSuccessEmail(requestId, user.Email, user.Name)
		}
	} else {
		logger.Error("[{}] error fetching user by identifier: {}", requestId, identifier)
		return nil, err
	}

	jwtToken, jwtError := tokenService.GetInstance().GenerateJwtToken(requestId, *user)

	if jwtError != nil {
		logger.Error("[{}]: Error generating auth token -> {}", requestId, jwtError)
		return nil, jwtError
	}

	authDao.UpdateTimestamp(requestId, identifier, authDao.TIMESTAMP_TYPE_LAST_LOGIN_TIME)

	return &model.OAuthCallbackResponse{
		AuthToken: jwtToken,
		UserId:    user.Id,
		Success:   true,
	}, nil
}

func registerUser(requestId string, profileInfo ProfileInfo) (*model.User, *model.ErrorResponse) {
	userToSave := createUserEntity(profileInfo)

	registeredUser, err := authDao.SaveUser(requestId, userToSave)

	if err != nil {
		return nil, err
	}

	return &model.User{
		Id:                  registeredUser.Id,
		Name:                registeredUser.Name,
		Email:               registeredUser.Email,
		Password:            registeredUser.Password,
		Scopes:              registeredUser.Scopes,
		ForgotPasswordToken: utils.GetStringOrNil(registeredUser.ForgotPasswordToken),
		LastLoginAt:         utils.GetInt64OrNil(registeredUser.LastLoginAt),
		PasswordChangedAt:   utils.GetInt64OrNil(registeredUser.LastPasswordChangedAt),
		IsDeleted:           registeredUser.IsDeleted,
	}, nil
}

func createUserEntity(profileInfo ProfileInfo) *entity.User {
	var loginType entity.LoginType
	var email string

	if profileInfo.Email != "" {
		loginType = entity.OAUTH_OTP
		email = profileInfo.Email
	} else {
		loginType = entity.OAUTH_ONLY
		email = profileInfo.Username
	}

	return &entity.User{
		Id:                strings.ReplaceAll(uuid.New().String(), "-", ""),
		Email:             email,
		ProfilePictureURL: &profileInfo.ProfilePicture,
		Name:              profileInfo.Name,
		UserLoginType:     loginType,
		Scopes:            "user",
	}
}

func getIdentifier(profileInfo ProfileInfo) string {
	if profileInfo.Email != "" {
		return profileInfo.Email
	}
	return profileInfo.Id
}

func getProfileInfo(reqId string, request model.OAuthCallbackRequest) (*ProfileInfo, *model.ErrorResponse) {
	oAuthProvider := request.Provider

	var profileInfo ProfileInfo

	switch oAuthProvider {
	case model.OAuthProvider(model.OAUTH_PROVIDER_GITHUB):
		pInfo, err := FetchGitHubProfileInfo(reqId, request)

		if err != nil {
			return nil, err
		}

		profileInfo = *pInfo

	case model.OAuthProvider(model.OAUTH_PROVIDER_GOOGLE):
		pInfo, err := FetchGoogleProfileInfo(reqId, request)

		if err != nil {
			return nil, err
		}

		profileInfo = *pInfo
	}

	logger.Debug("[{}] profile info fetched is: {}", reqId, profileInfo)

	return &profileInfo, nil
}
