package dao

import (
	"database/sql"
	"fmt"
	"time"

	MySQL "github.com/akgarg0472/urlshortener-auth-service/database"
	Models "github.com/akgarg0472/urlshortener-auth-service/model"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	utils "github.com/akgarg0472/urlshortener-auth-service/utils"
)

var (
	logger = Logger.GetLogger("authDao.go")
)

func GetUserByEmail(requestId string, email string) (*Models.User, *Models.ErrorResponse) {
	logger.Info("[{}]: Getting user by email -> {}", requestId, email)

	rows, err := doSelectQuery(requestId, SELECT_USER_BY_EMAIL_QUERY, email)

	if err != nil {
		logger.Error("[{}]: Error executing query: {}", requestId, err.Error())
		return nil, utils.InternalServerErrorResponse()
	}

	defer rows.Close()

	users := make([]Models.User, 0)

	for rows.Next() {
		var user Models.User
		var lastLoginAt, passwordChangedAt []uint8

		err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.Scopes, &user.Name)

		if err != nil {
			logger.Error("[{}]: Error parsing query result: {}", requestId, err.Error())
			return nil, utils.InternalServerErrorResponse()
		}

		lastLogin, err := time.Parse("2006-01-02 15:04:05", string(lastLoginAt))
		if err == nil {
			user.LastLoginAt = lastLogin
		}

		passwordChanged, err := time.Parse("2006-01-02 15:04:05", string(passwordChangedAt))
		if err == nil {
			user.PasswordChangedAt = passwordChanged
		}

		users = append(users, user)
	}

	if len(users) == 0 {
		logger.Error("[{}]: No user found by email: {}", requestId, email)
		return nil, utils.GetErrorResponse(fmt.Sprintf("No user found by email: %s", email), 404)
	} else if len(users) == 1 {
		logger.Debug("[{}]: user found with email", requestId)
		return &users[0], nil
	} else {
		logger.Error("[{}]: Multiple users found by email: {}", requestId, email)
		return nil, utils.InternalServerErrorResponse()
	}
}

func CheckIfUserExistsByEmail(requestId string, email string) (bool, *Models.ErrorResponse) {
	logger.Info("[{}]: Checking if user exists -> {}", requestId, email)

	fmt.Println(CHECK_USER_EXISTS_BY_EMAIL_QUERY)

	result, err := doSelectQuery(requestId, CHECK_USER_EXISTS_BY_EMAIL_QUERY, email)

	if err != nil {
		logger.Error("[{}]: Error executing query: {}", requestId, err.Error())
		return false, utils.InternalServerErrorResponse()
	}

	defer result.Close()

	var count int

	for result.Next() {
		err := result.Scan(&count)

		if err != nil {
			logger.Error("[{}]: Error parsing query result: {}", requestId, err.Error())
			return false, utils.InternalServerErrorResponse()
		}
	}

	logger.Debug("[{}]: CheckIfUserExistsByEmail Count Query Result -> {}", requestId, count)

	return count == 1, nil
}

func SaveUser(requestId string, signupRequest Models.SignupRequest) (*Models.User, *Models.ErrorResponse) {
	logger.Info("[{}]: Saving user into DB -> {}", requestId, signupRequest)

	db := MySQL.GetInstance(requestId, "authDao.go")

	if db == nil {
		logger.Error("[{}]: Error getting DB instance", requestId)
		return nil, utils.InternalServerErrorResponse()
	}

	preparedStatement, err := db.Prepare(INSERT_USER_QUERY)

	if err != nil {
		logger.Error("[{}]: Error preparing insert statement: {}", requestId, err.Error())
		return nil, utils.InternalServerErrorResponse()
	}

	defer preparedStatement.Close()

	user := createUserEntity(signupRequest)

	_, insertError := preparedStatement.Exec(user.Id, user.Email, user.Password, user.Scopes, user.Name)

	if insertError != nil {
		logger.Error("[{}]: Error saving user: {}", requestId, insertError.Error())
		return nil, utils.ParseMySQLErrorAndReturnErrorResponse(insertError)
	}

	return user, nil
}

// updates the value of token stored in DB for the identity (identity could be id or email of user)
func UpdateForgotPasswordToken(requestId string, identity string, token string) (bool, *Models.ErrorResponse) {
	logger.Info("[{}]: Updating user token info into DB", requestId)

	result, updateQueryError := doUpdateQuery(requestId, UPDATE_FORGOT_PASSWORD_TOKEN_QUERY, token, identity, identity)

	if updateQueryError != nil {
		return false, updateQueryError
	}

	rowsAffected, rowsAffectedError := result.RowsAffected()

	if rowsAffectedError == nil && rowsAffected == 1 {
		logger.Info("[{}]: Forgot password token updated", requestId)
		return true, nil
	}

	if rowsAffectedError != nil {
		logger.Error("[{}]: Error getting rows affected: {}", requestId, rowsAffectedError.Error())
	}

	return false, utils.InternalServerErrorResponse()
}

func GetForgotPasswordToken(requestId string, email string) (string, *Models.ErrorResponse) {
	logger.Info("[{}]: Fetching user forgot token info from DB", requestId)

	result, queryErr := doSelectQuery(requestId, GET_FORGOT_TOKEN_BY_EMAIL_QUERY, email)

	if queryErr != nil {
		return "", utils.InternalServerErrorResponse()
	}

	defer result.Close()

	var forgotPasswordToken sql.NullString

	for result.Next() {
		err := result.Scan(&forgotPasswordToken)

		if err != nil {
			logger.Error("[{}]: Error scanning result: {}", requestId, err.Error())
			return "", utils.InternalServerErrorResponse()
		}
	}

	if !forgotPasswordToken.Valid {
		logger.Error("[{}] invalid forgot password token fetched from DB", requestId)
		return "", utils.GetErrorResponse("Invalid Forgot Password Token. Request Rejected", 400)
	}

	return forgotPasswordToken.String, nil
}

func UpdatePassword(requestId string, identity string, newPassword string) (bool, *Models.ErrorResponse) {
	logger.Info("[{}]: Updating user password into DB", requestId)

	result, updateQueryError := doUpdateQuery(requestId, UPDATE_PASSWORD_QUERY, newPassword, identity, identity)

	if updateQueryError != nil {
		return false, updateQueryError
	}

	rowsAffected, rowsAffectedError := result.RowsAffected()

	if rowsAffectedError == nil && rowsAffected == 1 {
		logger.Info("[{}]: Forgot password token updated", requestId)
		return true, nil
	}

	if rowsAffectedError != nil {
		logger.Error("[{}]: Error getting rows affected: {}", requestId, rowsAffectedError.Error())
	}

	return false, utils.InternalServerErrorResponse()
}
