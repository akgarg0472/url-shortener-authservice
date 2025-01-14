package oauth_service

import (
	"fmt"
	"strings"

	enums "github.com/akgarg0472/urlshortener-auth-service/constants"
	entity2 "github.com/akgarg0472/urlshortener-auth-service/internal/entity"

	authDao "github.com/akgarg0472/urlshortener-auth-service/internal/dao/auth"
	oauthDao "github.com/akgarg0472/urlshortener-auth-service/internal/dao/oauth"
	kafka_service "github.com/akgarg0472/urlshortener-auth-service/internal/service/kafka"
	notificationService "github.com/akgarg0472/urlshortener-auth-service/internal/service/notification"
	tokenService "github.com/akgarg0472/urlshortener-auth-service/internal/service/token"
	"github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
	"github.com/google/uuid"
)

var (
	oAuthProvidersMapping map[string]model.OAuthProvider
	logger                = Logger.GetLogger("oAuthService.go")
)

type ProfileInfo struct {
	OAuthId        string
	Name           string
	ProfilePicture string
	Email          string
	OAuthProvider  string
}

func (p ProfileInfo) String() string {
	return fmt.Sprintf(`{"OAuthId":"%s", "Name":"%s", "ProfilePicture":"%s", "Email":"%s"}`, p.OAuthId, p.Name, p.ProfilePicture, p.Email)
}

type AccessTokenResponse struct {
	AccessToken string
	TokenType   string
}

func InitOAuthProviders() {
	logger.Info("Initializing oAuth providers")
	clients := oauthDao.FetchOAuthProviders()

	if oAuthProvidersMapping == nil {
		oAuthProvidersMapping = make(map[string]model.OAuthProvider)
	}

	for _, _client := range clients {
		client := model.OAuthProvider{
			Provider:    _client.Provider,
			ClientId:    _client.ClientID,
			BaseUrl:     _client.BaseUrl,
			RedirectURI: _client.RedirectURI,
			AccessType:  _client.AccessType,
			Scope:       _client.Scope,
		}

		oAuthProvidersMapping[client.Provider] = client
	}

	logger.Debug("Loaded oAuth providers: {}", oAuthProvidersMapping)
}

func GetOAuthProvider(query string) []model.OAuthProvider {
	var clients []model.OAuthProvider

	if query == "" {
		for _, client := range oAuthProvidersMapping {
			clients = append(clients, client)
		}

		return clients
	}

	providers := strings.Split(query, ",")

	for _, provider := range providers {
		if client, exists := oAuthProvidersMapping[provider]; exists {
			clients = append(clients, client)
		}
	}

	return clients
}

func ProcessCallbackRequest(
	requestId string,
	oAuthCallbackRequest model.OAuthCallbackRequest,
) (*model.OAuthCallbackResponse, *model.ErrorResponse) {
	logger.Info("[{}] oAuth callback request received: {}", requestId, oAuthCallbackRequest)
	var newUser bool
	profileInfo, err := getProfileInfo(requestId, oAuthCallbackRequest)

	if err != nil {
		logger.Error("[{}] error fetching oAuth profile details: {}", requestId, err)
		return nil, err
	}

	// checks if user is registered or not
	user, err := getExistingUser(requestId, *profileInfo)

	if user != nil {
		logger.Info("[{}] user is already registered", requestId)
		newUser = false
	} else if err.ErrorCode == 404 {
		logger.Info("[{}] user is not registered, going to register it", requestId)
		newUser = true
		registeredUser, err := registerUser(requestId, *profileInfo)

		if err != nil {
			logger.Error("[{}] failed to register OAuth user: {}", requestId, err)
			return nil, err
		}

		user = registeredUser

		if user.Email != "" {
			notificationService.SendSignupSuccessEmail(requestId, user.Email, user.Name)
		}

		// push user registered kafka event
		kafka_service.GetInstance().PushUserRegisteredEvent(requestId, user.Id)

	} else if err != nil && err.ErrorCode == 409 {
		logger.Error("[{}] user already exists for oauthId/email: {}, {}", requestId, profileInfo.OAuthId, profileInfo.Email)
		return nil, err
	} else {
		logger.Error("[{}] error fetching user: {}", requestId, profileInfo.OAuthId)
		return nil, err
	}

	jwtToken, jwtError := tokenService.GetInstance().GenerateJwtToken(requestId, *user)

	if jwtError != nil {
		logger.Error("[{}]: Error generating auth token -> {}", requestId, jwtError)
		return nil, jwtError
	}

	authDao.UpdateTimestamp(requestId, user.Id, authDao.TimestampTypeLastLoginTime)

	message := ""
	if newUser {
		message = "Welcome onboard: " + profileInfo.Name
	} else {
		message = "Welcome back: " + profileInfo.Name
	}

	return &model.OAuthCallbackResponse{
		AuthToken: jwtToken,
		UserId:    user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Success:   true,
		IsNewUser: newUser,
		Message:   message,
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
		Email:               utils.GetStringOrNil(registeredUser.Email),
		Scopes:              registeredUser.Scopes,
		ForgotPasswordToken: utils.GetStringOrNil(registeredUser.ForgotPasswordToken),
		LastLoginAt:         utils.GetInt64OrNil(registeredUser.LastLoginAt),
		PasswordChangedAt:   utils.GetInt64OrNil(registeredUser.LastPasswordChangedAt),
		IsDeleted:           registeredUser.IsDeleted,
	}, nil
}

func createUserEntity(profileInfo ProfileInfo) *entity2.User {
	var entityLoginType enums.UserEntityLoginType
	var email *string

	if profileInfo.Email != "" {
		entityLoginType = enums.UserEntityLoginTypeOauthAndOtp
		email = &profileInfo.Email
	} else {
		entityLoginType = enums.UserEntityLoginTypeOauthOnly
		email = nil
	}

	return &entity2.User{
		Id:                strings.ReplaceAll(uuid.New().String(), "-", ""),
		OAuthId:           &profileInfo.OAuthId,
		Email:             email,
		ProfilePictureURL: &profileInfo.ProfilePicture,
		Name:              profileInfo.Name,
		UserLoginType:     entityLoginType,
		Scopes:            "user",
		OAuthProvider:     &profileInfo.OAuthProvider,
	}
}

func getProfileInfo(reqId string, request model.OAuthCallbackRequest) (*ProfileInfo, *model.ErrorResponse) {
	oAuthProvider := request.Provider

	var profileInfo ProfileInfo

	switch oAuthProvider {
	case enums.OauthProviderGithub:
		pInfo, err := FetchGitHubProfileInfo(reqId, request)
		if err != nil {
			return nil, err
		}
		profileInfo = *pInfo
		profileInfo.OAuthProvider = string(enums.OauthProviderGithub)

	case enums.OauthProviderGoogle:
		pInfo, err := FetchGoogleProfileInfo(reqId, request)
		if err != nil {
			return nil, err
		}
		profileInfo = *pInfo
		profileInfo.OAuthProvider = string(enums.OauthProviderGoogle)
	}

	logger.Debug("[{}] profile info fetched is: {}", reqId, profileInfo)

	return &profileInfo, nil
}

func getExistingUser(requestId string, profileInfo ProfileInfo) (*model.User, *model.ErrorResponse) {
	user, err := authDao.GetUserByOAuthId(requestId, profileInfo.OAuthId)

	if err != nil {
		if err.ErrorCode == 404 {
			if profileInfo.Email == "" {
				return nil, &model.ErrorResponse{
					ErrorCode: 404,
				}
			}

			userExistsByEmail, emailError := authDao.CheckIfUserExistsByEmail(requestId, profileInfo.Email)

			if emailError != nil {
				return nil, emailError
			}

			if userExistsByEmail {
				return nil, &model.ErrorResponse{
					Message:   "An account exists by same email: " + profileInfo.Email,
					ErrorCode: 409,
				}
			}
		}

		return nil, err
	}

	if user != nil {
		return user, nil
	} else {
		return nil, nil
	}
}
