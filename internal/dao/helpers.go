package dao

import (
	"database/sql"
	"fmt"
	"strings"

	MySQL "github.com/akgarg0472/urlshortener-auth-service/database"
	Models "github.com/akgarg0472/urlshortener-auth-service/model"
	utils "github.com/akgarg0472/urlshortener-auth-service/utils"
	"github.com/google/uuid"
)

func doSelectQuery(requestId string, query string, params ...interface{}) (*sql.Rows, error) {
	logger.Debug("[{}]: Executing query -> {}", requestId, query)

	var db = MySQL.GetInstance(requestId, "doSelectQuery")

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

func doUpdateQuery(requestId string, query string, params ...interface{}) (sql.Result, *Models.ErrorResponse) {
	db := MySQL.GetInstance(requestId, "authDao.go")

	if db == nil {
		logger.Error("[{}]: Error getting DB instance", requestId)
		return nil, utils.InternalServerErrorResponse()
	}

	preparedStatement, err := db.Prepare(query)

	if err != nil {
		logger.Error("[{}]: Error preparing statement: {}", requestId, err.Error())
		return nil, utils.InternalServerErrorResponse()
	}

	defer preparedStatement.Close()

	result, err := preparedStatement.Exec(params...)

	if err != nil {
		logger.Error("[{}]: Error performing update query: {}", requestId, err.Error())
		return nil, utils.ParseMySQLErrorAndReturnErrorResponse(err)
	}

	return result, nil
}

func createUserEntity(request Models.SignupRequest) *Models.User {
	return &Models.User{
		Id:       strings.ReplaceAll(uuid.New().String(), "-", ""),
		Email:    request.Email,
		Password: request.Password,
		Name:     request.Name,
		Scopes:   "user",
	}
}
