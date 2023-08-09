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

	rows, err := doSelectQuery(requestId, "SELECT * FROM users WHERE email = ?", email)

	if err != nil {
		logger.Error("[{}]: Error executing statement: {}", requestId, err.Error())
		return nil, utils.InternalServerErrorResponse()
	}

	defer rows.Close()

	users := make([]Models.User, 0)

	for rows.Next() {
		var user Models.User
		var createdAt, updatedAt []uint8

		err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.Scopes, &createdAt, &updatedAt)

		if err != nil {
			logger.Error("[{}]: Error scanning result: {}", requestId, err.Error())
			return nil, utils.InternalServerErrorResponse()
		}

		createdAtTime, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err == nil {
			user.CreatedAt = createdAtTime
		}

		updatedAtTime, err := time.Parse("2006-01-02 15:04:05", string(updatedAt))
		if err == nil {
			user.UpdatedAt = updatedAtTime
		}

		users = append(users, user)
	}

	logger.Debug("[{}]: UserByEmail Query Result -> {}", requestId, users)

	if len(users) == 0 {
		logger.Error("[{}]: No user found by email: {}", requestId, email)
		return nil, utils.GetErrorResponse(fmt.Sprintf("No user found by email: %s", email), 404)
	} else if len(users) == 1 {
		logger.Debug("[{}]: Returning user -> {}", requestId, users[0])
		return &users[0], nil
	} else {
		logger.Error("[{}]: Multiple users found by email: {}", requestId, email)
		return nil, utils.InternalServerErrorResponse()
	}
}

func CheckIfUserExistsByEmail(requestId string, email string) (bool, *Models.ErrorResponse) {
	logger.Info("[{}]: Checking if user exists -> {}", requestId, email)

	result, err := doSelectQuery(requestId, "SELECT count(*) FROM users WHERE email = ?", email)

	if err != nil {
		logger.Error("[{}]: Error executing statement: {}", requestId, err.Error())
		return false, utils.InternalServerErrorResponse()
	}

	defer result.Close()

	var count int

	for result.Next() {
		err := result.Scan(&count)

		if err != nil {
			logger.Error("[{}]: Error scanning result: {}", requestId, err.Error())
			return false, utils.InternalServerErrorResponse()
		}
	}

	logger.Debug("[{}]: CheckIfUserExistsByEmail Count Query Result -> {}", requestId, count)

	return count > 0, nil
}

func SaveUser(requestId string, signupRequest Models.SignupRequest) (bool, *Models.ErrorResponse) {
	logger.Info("[{}]: Saving user into DB -> {}", requestId, signupRequest)

	var db = MySQL.GetInstance(requestId, "authDao.go")

	if db == nil {
		logger.Error("[{}]: Error getting DB instance", requestId)
		return false, utils.InternalServerErrorResponse()
	}

	preparedStatement, err := db.Prepare("INSERT INTO users (id, email, password, scopes) VALUES (?, ?, ?, ?)")

	if err != nil {
		logger.Error("[{}]: Error preparing statement: {}", requestId, err.Error())
		return false, utils.InternalServerErrorResponse()
	}

	defer preparedStatement.Close()

	result, err := preparedStatement.Exec(signupRequest.UserId, signupRequest.Email, signupRequest.Password, "user")

	if err != nil {
		logger.Error("[{}]: Error saving user: {}", requestId, err.Error())
		return false, utils.ParseMySQLErrorAndReturnErrorResponse(err)
	}

	logger.Debug("[{}]: Insert Result -> {}", requestId, result)

	rowsAffected, rowsAffectedError := result.RowsAffected()

	logger.Debug("[{}]: Rows affected -> {}", requestId, rowsAffected)

	if rowsAffectedError != nil {
		logger.Error("[{}]: Error getting rows affected: {}", requestId, rowsAffectedError.Error())
		return false, utils.InternalServerErrorResponse()
	}

	logger.Info("[{}]: User saved successfully -> {}", requestId, rowsAffected == 1)

	return rowsAffected == 1, nil
}

func doSelectQuery(requestId string, query string, params ...interface{}) (*sql.Rows, error) {
	logger.Debug("[{}]: Executing query -> {}", requestId, query)

	var db = MySQL.GetInstance(requestId, "authDao.go")

	if db == nil {
		logger.Error("[{}]: Error getting DB instance", requestId)
		return nil, fmt.Errorf("error getting DB instance")
	}

	preparedStatement, err := db.Prepare(query)

	if err != nil {
		logger.Error("[{}]: Error preparing statement: {}", requestId, err.Error())
		return nil, err
	}

	defer preparedStatement.Close()

	rows, err := preparedStatement.Query(params...)

	if err != nil {
		logger.Error("[{}]: Error executing statement: {}", requestId, err.Error())
		return nil, err
	}

	return rows, nil
}
