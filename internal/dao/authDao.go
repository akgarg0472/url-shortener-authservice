package dao

import (
	"fmt"
	"time"

	MySQL "github.com/akgarg0472/urlshortener-auth-service/database"
	"github.com/akgarg0472/urlshortener-auth-service/internal/dao/entity"
	"github.com/akgarg0472/urlshortener-auth-service/model"
	Models "github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	utils "github.com/akgarg0472/urlshortener-auth-service/utils"
	"gorm.io/gorm"
)

var (
	logger = Logger.GetLogger("authDao.go")
)

type TimestampType string

const (
	TIMESTAMP_TYPE_LAST_LOGIN_TIME TimestampType = "last_login_at"
)

func GetUserByEmail(requestId string, email string) (*Models.User, *Models.ErrorResponse) {
	logger.Info("[{}]: Getting user by email -> {}", requestId, email)

	db := MySQL.GetInstance(requestId, "GetUserByEmail")

	if db == nil {
		logger.Error("[{}]: Error getting DB instance", requestId)
		return nil, utils.InternalServerErrorResponse()
	}

	var dbUser entity.User

	result := db.First(&dbUser, "email = ?", email)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			logger.Error("[{}]: No user found with email: {}", requestId, email)
			return nil, utils.GetErrorResponse(fmt.Sprintf("No user found by email: %s", email), 404)
		} else {
			logger.Error("[{}]: Error querying user=: {}", requestId, result.Error)
		}

		return nil, utils.InternalServerErrorResponse()
	}

	user := model.User{
		Id:                  dbUser.Id,
		Name:                dbUser.Name,
		Email:               dbUser.Email,
		Password:            dbUser.Password,
		Scopes:              dbUser.Scopes,
		ForgotPasswordToken: getStringOrNil(dbUser.ForgotPasswordToken),
		LastLoginAt:         getInt64OrNil(dbUser.LastLoginAt),
		PasswordChangedAt:   getInt64OrNil(dbUser.LastPasswordChangedAt),
		IsDeleted:           dbUser.IsDeleted,
	}

	logger.Info("[{}] Fetched user: {}", requestId, user)

	return &user, nil
}

func CheckIfUserExistsByEmail(requestId string, email string) (bool, *Models.ErrorResponse) {
	logger.Info("[{}]: Checking if user exists -> {}", requestId, email)

	db := MySQL.GetInstance(requestId, "CheckIfUserExistsByEmail")

	if db == nil {
		logger.Error("[{}]: Error getting DB instance", requestId)
		return false, utils.InternalServerErrorResponse()
	}

	var count int64
	result := db.Model(&entity.User{}).Where("email = ?", email).Count(&count)

	if result.Error != nil {
		logger.Error("[{}] error checking for user existence by email", requestId, result.Error)
		return false, utils.InternalServerErrorResponse()
	}

	logger.Info("[{}]: CheckIfUserExistsByEmail Count Query Result -> {}", requestId, count)
	return count == 1, nil
}

func SaveUser(requestId string, user *entity.User) (*entity.User, *Models.ErrorResponse) {
	logger.Info("[{}]: Saving user into DB -> {}", requestId, user.Email)

	db := MySQL.GetInstance(requestId, "SaveUser")

	if db == nil {
		logger.Error("[{}]: Error getting DB instance", requestId)
		return nil, utils.InternalServerErrorResponse()
	}

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
		logger.Error("[{}]: Error getting DB instance", requestId)
		return false, utils.InternalServerErrorResponse()
	}

	var affactedRows int64
	result := db.Model(&entity.User{}).Where("id = ? or email = ?", identity, identity).UpdateColumn("ForgotPasswordToken", token).Count(&affactedRows)

	if result.Error != nil {
		logger.Error("[{}]: Error updating forgot password token: {}", requestId, result.Error)
		return false, utils.InternalServerErrorResponse()
	}

	if affactedRows == 1 {
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
		logger.Error("[{}]: Error getting DB instance", requestId)
		return false, utils.InternalServerErrorResponse()
	}

	timestamp := time.Now().UnixMilli()
	var affactedRows int64

	result := db.Model(&entity.User{}).Where("id = ? or email = ?", identity, identity).UpdateColumns(map[string]interface{}{
		"password":                 newPassword,
		"forgot_password_token":    "", // Empty string for the token
		"last_password_changed_at": timestamp,
	}).Count(&affactedRows)

	if result.Error != nil {
		logger.Error("[{}]: Error updating password: {}", requestId, result.Error)
		return false, utils.InternalServerErrorResponse()
	}

	if affactedRows == 1 {
		logger.Info("[{}]: Password updated successfully", requestId)
		return true, nil
	}

	return false, utils.InternalServerErrorResponse()
}

func UpdateTimestamp(requestId string, email string, timestampType TimestampType) {
	timestamp := time.Now().UnixMilli()

	logger.Debug("[{}]: Updating {} into DB: {}", requestId, timestampType, timestamp)

	db := MySQL.GetInstance(requestId, "UpdateTimestamp")

	if db == nil {
		logger.Error("[{}]: Error getting DB instance", requestId)
		return
	}

	var affactedRows int64

	result := db.Model(&entity.User{}).Where("email = ?", email).UpdateColumn(string(timestampType), timestamp).Count(&affactedRows)

	if result.Error != nil {
		logger.Error("[{}]: Error updating {}: {}", requestId, timestampType, result.Error)
		return
	}

	if affactedRows == 1 {
		logger.Info("[{}] {} updated successfully", requestId, timestampType)
		return
	}
}
