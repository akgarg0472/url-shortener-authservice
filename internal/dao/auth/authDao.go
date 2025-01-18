package auth_dao

import (
	"errors"
	"fmt"
	"time"

	"github.com/akgarg0472/urlshortener-auth-service/internal/entity"

	MySQL "github.com/akgarg0472/urlshortener-auth-service/database"
	"github.com/akgarg0472/urlshortener-auth-service/model"
	Models "github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
	"gorm.io/gorm"
)

var (
	logger = Logger.GetLogger("authDao.go")
)

type TimestampType string

const (
	TimestampTypeLastLoginTime TimestampType = "LastLoginAt"
)

func logErrorGettingDBInstance(requestId string) {
	logger.Error("[{}]: Error getting DB instance", requestId)
}

func GetUserByEmail(requestId string, identity string) (*Models.User, *Models.ErrorResponse) {
	logger.Info("[{}]: Getting user by email -> {}", requestId, identity)

	db := MySQL.GetInstance(requestId, "GetUserByEmail")

	if db == nil {
		logErrorGettingDBInstance(requestId)
		return nil, utils.InternalServerErrorResponse()
	}

	var dbUser entity.User

	result := db.First(&dbUser, "email = ?", identity)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Info("[{}]: No user found with email: {}", requestId, identity)
			return nil, utils.GetErrorResponse("email not registered", 404)
		} else {
			logger.Error("[{}]: Error querying user=: {}", requestId, result.Error)
		}

		return nil, utils.InternalServerErrorResponse()
	}

	user := model.User{
		Id:                  dbUser.Id,
		Name:                dbUser.Name,
		Email:               utils.GetStringOrNil(dbUser.Email),
		Password:            utils.GetStringOrNil(dbUser.Password),
		Scopes:              dbUser.Scopes,
		ForgotPasswordToken: utils.GetStringOrNil(dbUser.ForgotPasswordToken),
		LastLoginAt:         utils.GetInt64OrNil(dbUser.LastLoginAt),
		PasswordChangedAt:   utils.GetInt64OrNil(dbUser.LastPasswordChangedAt),
		IsDeleted:           dbUser.IsDeleted,
		OAuthId:             utils.GetStringOrNil(dbUser.OAuthId),
		LoginType:           dbUser.UserLoginType,
		OAuthProvider:       utils.GetStringOrNil(dbUser.OAuthProvider),
	}

	logger.Debug("[{}] Fetched user: {}", requestId, user)

	return &user, nil
}

func GetUserById(requestId string, identity string) (*Models.User, *Models.ErrorResponse) {
	logger.Info("[{}]: Getting user by id -> {}", requestId, identity)

	db := MySQL.GetInstance(requestId, "GetUserById")

	if db == nil {
		logErrorGettingDBInstance(requestId)
		return nil, utils.InternalServerErrorResponse()
	}

	var dbUser entity.User

	result := db.First(&dbUser, "id = ?", identity)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Info("[{}]: No user found with id: {}", requestId, identity)
			return nil, utils.GetErrorResponse("User not found with id", 404)
		} else {
			logger.Error("[{}]: Error querying user=: {}", requestId, result.Error)
		}

		return nil, utils.InternalServerErrorResponse()
	}

	user := model.User{
		Id:                  dbUser.Id,
		Name:                dbUser.Name,
		Email:               utils.GetStringOrNil(dbUser.Email),
		Password:            utils.GetStringOrNil(dbUser.Password),
		Scopes:              dbUser.Scopes,
		ForgotPasswordToken: utils.GetStringOrNil(dbUser.ForgotPasswordToken),
		LastLoginAt:         utils.GetInt64OrNil(dbUser.LastLoginAt),
		PasswordChangedAt:   utils.GetInt64OrNil(dbUser.LastPasswordChangedAt),
		IsDeleted:           dbUser.IsDeleted,
		OAuthId:             utils.GetStringOrNil(dbUser.OAuthId),
		LoginType:           dbUser.UserLoginType,
		OAuthProvider:       utils.GetStringOrNil(dbUser.OAuthProvider),
	}

	logger.Debug("[{}] Fetched user: {}", requestId, user)

	return &user, nil
}

func GetUserByOAuthId(requestId string, oAuthId string) (*Models.User, *Models.ErrorResponse) {
	logger.Info("[{}]: Getting user by oAuthId -> {}", requestId, oAuthId)

	db := MySQL.GetInstance(requestId, "GetUserByOAuthId")

	if db == nil {
		logErrorGettingDBInstance(requestId)
		return nil, utils.InternalServerErrorResponse()
	}

	var dbUser entity.User

	result := db.First(&dbUser, "oauth_id = ?", oAuthId)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Error("[{}]: No user found with oAuthId: {}", requestId, oAuthId)
			return nil, utils.GetErrorResponse(fmt.Sprintf("No user found by oAuthId: %s", oAuthId), 404)
		} else {
			logger.Error("[{}]: Error querying user=: {}", requestId, result.Error)
		}

		return nil, utils.InternalServerErrorResponse()
	}

	user := model.User{
		Id:                  dbUser.Id,
		Name:                dbUser.Name,
		Email:               utils.GetStringOrNil(dbUser.Email),
		Password:            utils.GetStringOrNil(dbUser.Password),
		Scopes:              dbUser.Scopes,
		ForgotPasswordToken: utils.GetStringOrNil(dbUser.ForgotPasswordToken),
		LastLoginAt:         utils.GetInt64OrNil(dbUser.LastLoginAt),
		PasswordChangedAt:   utils.GetInt64OrNil(dbUser.LastPasswordChangedAt),
		IsDeleted:           dbUser.IsDeleted,
		OAuthId:             *dbUser.OAuthId,
		OAuthProvider:       *dbUser.OAuthProvider,
	}

	logger.Info("[{}] Fetched user: {}", requestId, user)

	return &user, nil
}

func CheckIfUserExistsByEmail(requestId string, email string) (bool, *Models.ErrorResponse) {
	logger.Info("[{}]: Checking if user exists by email: {}", requestId, email)

	db := MySQL.GetInstance(requestId, "CheckIfUserExistsByEmail")

	if db == nil {
		logErrorGettingDBInstance(requestId)
		return false, utils.InternalServerErrorResponse()
	}

	var count int64
	result := db.Model(&entity.User{}).Where("email = ?", email).Count(&count)

	if result.Error != nil {
		logger.Error("[{}] error checking for user existence by email", requestId, result.Error)
		return false, utils.InternalServerErrorResponse()
	}

	logger.Info("[{}]: CheckIfUserExistsByEmail Count Query Result: {}", requestId, count)
	return count == 1, nil
}

func SaveUser(requestId string, user *entity.User) (*entity.User, *Models.ErrorResponse) {
	logger.Info("[{}]: Saving user into DB with id: {} and email: {}", requestId, user.Id, utils.GetStringOrNil(user.Email))

	db := MySQL.GetInstance(requestId, "SaveUser")

	if db == nil {
		logErrorGettingDBInstance(requestId)
		return nil, utils.InternalServerErrorResponse()
	}

	user.CreatedAt = time.Now().UnixMilli()
	user.UpdatedAt = time.Now().UnixMilli()

	result := db.Create(user)

	if result.Error != nil {
		logger.Error("[{}]: Error saving user: {}", requestId, result.Error)
		return nil, utils.InternalServerErrorResponse()
	}

	logger.Info("[{}] user created successfully", requestId)

	return user, nil
}

func UpdateForgotPasswordToken(requestId string, identity string, token string) (bool, *Models.ErrorResponse) {
	logger.Info("[{}]: Updating forgot password token: {}", requestId, identity)

	db := MySQL.GetInstance(requestId, "UpdateForgotPasswordToken")

	if db == nil {
		logErrorGettingDBInstance(requestId)
		return false, utils.InternalServerErrorResponse()
	}

	var affectedRows int64
	timestamp := time.Now().UnixMilli()

	result := db.Model(&entity.User{}).Where("id = ? or email = ?", identity, identity).UpdateColumns(
		map[string]interface{}{
			"ForgotPasswordToken": token,
			"UpdatedAt":           timestamp,
		},
	).Count(&affectedRows)

	if result.Error != nil {
		logger.Error("[{}]: Error updating forgot password token: {}", requestId, result.Error)
		return false, utils.InternalServerErrorResponse()
	}

	if affectedRows == 1 {
		logger.Info("[{}]: Forgot password token updated", requestId)
		return true, nil
	}

	return false, utils.InternalServerErrorResponse()
}

func GetForgotPasswordToken(requestId string, email string) (string, *Models.ErrorResponse) {
	logger.Info("[{}]: Fetching user forgot token info from DB", requestId)

	user, err := GetUserByEmail(requestId, email)

	if err != nil {
		return "", err
	}

	if user.ForgotPasswordToken == "" {
		logger.Error("[{}] invalid forgot password token fetched from DB", requestId)
		return "", utils.GetErrorResponse("Invalid Forgot Password Token. Request Rejected", 400)
	}

	return user.ForgotPasswordToken, nil
}

func UpdatePassword(requestId string, identity string, newPassword string) (bool, *Models.ErrorResponse) {
	logger.Info("[{}]: Updating user password into DB", requestId)

	db := MySQL.GetInstance(requestId, "UpdatePassword")

	if db == nil {
		logErrorGettingDBInstance(requestId)
		return false, utils.InternalServerErrorResponse()
	}

	timestamp := time.Now().UnixMilli()
	var affectedRows int64

	result := db.Model(&entity.User{}).Where("id = ? or email = ?", identity, identity).UpdateColumns(map[string]interface{}{
		"password":              newPassword,
		"ForgotPasswordToken":   "", // Empty string for the token
		"LastPasswordChangedAt": timestamp,
		"UpdatedAt":             timestamp,
	}).Count(&affectedRows)

	if result.Error != nil {
		logger.Error("[{}]: Error updating password: {}", requestId, result.Error)
		return false, utils.InternalServerErrorResponse()
	}

	if affectedRows == 1 {
		logger.Info("[{}]: Password updated successfully", requestId)
		return true, nil
	}

	return false, utils.InternalServerErrorResponse()
}

func UpdateTimestamp(requestId string, identity string, timestampType TimestampType) {
	timestamp := time.Now().UnixMilli()

	logger.Debug("[{}]: Updating {} into DB: {}", requestId, timestampType, timestamp)

	db := MySQL.GetInstance(requestId, "UpdateTimestamp")

	if db == nil {
		logErrorGettingDBInstance(requestId)
		return
	}

	var affectedRows int64

	result := db.Model(&entity.User{}).Where("email = ? or id = ?", identity, identity).UpdateColumns(map[string]interface{}{
		string(timestampType): timestamp,
	}).Count(&affectedRows)

	if result.Error != nil {
		logger.Error("[{}]: Error updating {}: {}", requestId, timestampType, result.Error)
		return
	}

	if affectedRows == 1 {
		logger.Info("[{}] {} updated successfully", requestId, timestampType)
		return
	}
}
