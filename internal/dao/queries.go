package dao

import (
	"fmt"

	"github.com/akgarg0472/urlshortener-auth-service/utils"
)

var (
	userTableName                      = ""
	SELECT_USER_BY_EMAIL_QUERY         = ""
	CHECK_USER_EXISTS_BY_EMAIL_QUERY   = ""
	INSERT_USER_QUERY                  = ""
	UPDATE_FORGOT_PASSWORD_TOKEN_QUERY = ""
	GET_FORGOT_TOKEN_BY_EMAIL_QUERY    = ""
	UPDATE_PASSWORD_QUERY              = ""
)

func InitQueries() {
	userTableName = utils.GetEnvVariable("MYSQL_USERS_TABLE_NAME", "")

	SELECT_USER_BY_EMAIL_QUERY = fmt.Sprintf("SELECT id, email, password, scopes, name FROM %s WHERE email = ?", userTableName)
	CHECK_USER_EXISTS_BY_EMAIL_QUERY = fmt.Sprintf("SELECT count(1) FROM %s WHERE email = ?", userTableName)
	INSERT_USER_QUERY = fmt.Sprintf("INSERT INTO %s (id, email, password, scopes, name) VALUES (?, ?, ?, ?, ?)", userTableName)
	UPDATE_FORGOT_PASSWORD_TOKEN_QUERY = fmt.Sprintf("UPDATE %s SET forgot_password_token=? WHERE id=? OR email=?", userTableName)
	GET_FORGOT_TOKEN_BY_EMAIL_QUERY = fmt.Sprintf("SELECT forgot_password_token FROM %s WHERE email=?", userTableName)
	UPDATE_PASSWORD_QUERY = fmt.Sprintf("UPDATE %s SET password=? WHERE id=? OR email=?", userTableName)
}
