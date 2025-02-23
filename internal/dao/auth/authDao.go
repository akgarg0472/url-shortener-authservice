package auth_dao

import (
	"errors"
	"fmt"
	"time"

	"github.com/akgarg0472/urlshortener-auth-service/constants"
	"github.com/akgarg0472/urlshortener-auth-service/internal/entity"
	"github.com/akgarg0472/urlshortener-auth-service/internal/logger"
	"go.uber.org/zap"

	MySQL "github.com/akgarg0472/urlshortener-auth-service/database"
	Models "github.com/akgarg0472/urlshortener-auth-service/model"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
	"gorm.io/gorm"
)

type TimestampType string

const (
	TimestampTypeLastLoginTime TimestampType = "LastLoginAt"
)

func logErrorGettingDBInstance(requestId string) {
	if logger.IsErrorEnabled() {
		logger.Error("Error getting DB instance",
			zap.String(constants.RequestIdLogKey, requestId),
		)
	}
}

func GetUserByEmail(requestId string, identity string) (*Models.User, *Models.ErrorResponse) {
	if logger.IsInfoEnabled() {
		logger.Info("Getting user by email",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.String("email", identity),
		)
	}

	db := MySQL.GetInstance(requestId, "GetUserByEmail")

	if db == nil {
		logErrorGettingDBInstance(requestId)
		return nil, utils.InternalServerErrorResponse()
	}

	var dbUser entity.User

	result := db.First(&dbUser, "email = ?", identity)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			if logger.IsInfoEnabled() {
				logger.Info("No user found with email",
					zap.String(constants.RequestIdLogKey, requestId),
					zap.String("email", identity),
				)
			}
			return nil, utils.GetErrorResponse("email not registered", 404)
		} else {
			if logger.IsErrorEnabled() {
				logger.Error("Error querying user",
					zap.String(constants.RequestIdLogKey, requestId),
					zap.Error(result.Error),
				)
			}
		}
		return nil, utils.InternalServerErrorResponse()
	}

	user := Models.User{
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

	if logger.IsDebugEnabled() {
		logger.Debug("Fetched user",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.Any("user", user),
		)
	}

	return &user, nil
}

func GetUserById(requestId string, identity string) (*Models.User, *Models.ErrorResponse) {
	if logger.IsInfoEnabled() {
		logger.Info("Getting user by id",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.String("user_id", identity),
		)
	}

	db := MySQL.GetInstance(requestId, "GetUserById")

	if db == nil {
		logErrorGettingDBInstance(requestId)
		return nil, utils.InternalServerErrorResponse()
	}

	var dbUser entity.User

	result := db.First(&dbUser, "id = ?", identity)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			if logger.IsInfoEnabled() {
				logger.Info("No user found with id",
					zap.String(constants.RequestIdLogKey, requestId),
					zap.String("user_id", identity),
				)
			}
			return nil, utils.GetErrorResponse("User not found with id", 404)
		} else {
			if logger.IsErrorEnabled() {
				logger.Error("Error querying user",
					zap.String(constants.RequestIdLogKey, requestId),
					zap.Error(result.Error),
				)
			}
		}

		return nil, utils.InternalServerErrorResponse()
	}

	user := Models.User{
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

	if logger.IsDebugEnabled() {
		logger.Debug("Fetched user",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.Any("user", user),
		)
	}

	return &user, nil
}

func GetUserByOAuthId(requestId string, oAuthId string) (*Models.User, *Models.ErrorResponse) {
	if logger.IsInfoEnabled() {
		logger.Info("Getting user by oAuthId",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.String("oAuthId", oAuthId),
		)
	}

	db := MySQL.GetInstance(requestId, "GetUserByOAuthId")

	if db == nil {
		logErrorGettingDBInstance(requestId)
		return nil, utils.InternalServerErrorResponse()
	}

	var dbUser entity.User

	result := db.First(&dbUser, "oauth_id = ?", oAuthId)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			if logger.IsErrorEnabled() {
				logger.Error("No user found with oAuthId",
					zap.String(constants.RequestIdLogKey, requestId),
					zap.String("oAuthId", oAuthId),
				)
			}
			return nil, utils.GetErrorResponse(fmt.Sprintf("No user found by oAuthId: %s", oAuthId), 404)
		} else {
			if logger.IsErrorEnabled() {
				logger.Error("Error querying user",
					zap.String(constants.RequestIdLogKey, requestId),
					zap.Error(result.Error),
				)
			}
		}

		return nil, utils.InternalServerErrorResponse()
	}

	user := Models.User{
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
		LoginType:           dbUser.UserLoginType,
	}

	if logger.IsDebugEnabled() {
		logger.Debug("Fetched user",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.Any("user", user),
		)
	}

	return &user, nil
}

func CheckIfUserExistsByEmail(requestId string, email string) (bool, *Models.ErrorResponse) {
	if logger.IsInfoEnabled() {
		logger.Info("Checking if user exists by email",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.String("email", email),
		)
	}

	db := MySQL.GetInstance(requestId, "CheckIfUserExistsByEmail")

	if db == nil {
		logErrorGettingDBInstance(requestId)
		return false, utils.InternalServerErrorResponse()
	}

	var count int64
	result := db.Model(&entity.User{}).Where("email = ?", email).Count(&count)

	if result.Error != nil {
		if logger.IsErrorEnabled() {
			logger.Error("Error checking for user existence by email",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Error(result.Error),
			)
		}
		return false, utils.InternalServerErrorResponse()
	}

	if logger.IsInfoEnabled() {
		logger.Info("CheckIfUserExistsByEmail Count Query Result",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.Int64("count", count),
		)
	}

	return count == 1, nil
}

func SaveUser(requestId string, user *entity.User) (*entity.User, *Models.ErrorResponse) {
	if logger.IsInfoEnabled() {
		logger.Info("Saving user into DB",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.String("user_id", user.Id),
			zap.String("email", utils.GetStringOrNil(user.Email)),
		)
	}

	db := MySQL.GetInstance(requestId, "SaveUser")

	if db == nil {
		logErrorGettingDBInstance(requestId)
		return nil, utils.InternalServerErrorResponse()
	}

	user.CreatedAt = time.Now().UnixMilli()
	user.UpdatedAt = time.Now().UnixMilli()

	result := db.Create(user)

	if result.Error != nil {
		if logger.IsErrorEnabled() {
			logger.Error("Error saving user",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Error(result.Error),
			)
		}
		return nil, utils.InternalServerErrorResponse()
	}

	if logger.IsInfoEnabled() {
		logger.Info("User created successfully",
			zap.String(constants.RequestIdLogKey, requestId),
		)
	}

	return user, nil
}

func UpdateForgotPasswordToken(requestId string, identity string, token string) (bool, *Models.ErrorResponse) {
	if logger.IsInfoEnabled() {
		logger.Info("Updating forgot password token",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.String("identity", identity),
		)
	}

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
		if logger.IsErrorEnabled() {
			logger.Error("Error updating forgot password token",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Error(result.Error),
			)
		}
		return false, utils.InternalServerErrorResponse()
	}

	if affectedRows == 1 {
		if logger.IsInfoEnabled() {
			logger.Info("Forgot password token updated",
				zap.String(constants.RequestIdLogKey, requestId),
			)
		}
		return true, nil
	}

	return false, utils.InternalServerErrorResponse()
}

func GetForgotPasswordToken(requestId string, email string) (string, *Models.ErrorResponse) {
	if logger.IsInfoEnabled() {
		logger.Info("Fetching user forgot token info from DB",
			zap.String(constants.RequestIdLogKey, requestId),
		)
	}

	user, err := GetUserByEmail(requestId, email)

	if err != nil {
		return "", err
	}

	if user.ForgotPasswordToken == "" {
		if logger.IsErrorEnabled() {
			logger.Error("Invalid forgot password token fetched from DB",
				zap.String(constants.RequestIdLogKey, requestId),
			)
		}
		return "", utils.GetErrorResponse("Invalid Forgot Password Token", 400)
	}

	return user.ForgotPasswordToken, nil
}

func UpdatePassword(requestId string, identity string, newPassword string) (bool, *Models.ErrorResponse) {
	if logger.IsInfoEnabled() {
		logger.Info("Updating user password into DB",
			zap.String(constants.RequestIdLogKey, requestId),
		)
	}

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
		if logger.IsErrorEnabled() {
			logger.Error("Error updating password",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Error(result.Error),
			)
		}
		return false, utils.InternalServerErrorResponse()
	}

	if affectedRows == 1 {
		if logger.IsInfoEnabled() {
			logger.Info(
				"Password updated successfully",
				zap.String(constants.RequestIdLogKey, requestId),
			)
		}
		return true, nil
	}

	return false, utils.InternalServerErrorResponse()
}

func UpdateTimestamp(requestId string, identity string, timestampType TimestampType) {
	timestamp := time.Now().UnixMilli()

	if logger.IsDebugEnabled() {
		logger.Debug(
			"Updating timestamp into DB",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.Any("timestampType", timestampType),
			zap.Any("timestamp", timestamp),
		)
	}

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
		if logger.IsErrorEnabled() {
			logger.Error(
				"Error updating timestamp",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Any("timestampType", timestampType),
				zap.Error(result.Error),
			)
		}
		return
	}

	if affectedRows == 1 {
		if logger.IsInfoEnabled() {
			logger.Info(
				"Timestamp updated successfully",
				zap.String(constants.RequestIdLogKey, requestId),
				zap.Any("timestampType", timestampType),
			)
		}
		return
	}
}
