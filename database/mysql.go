package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"

	_ "github.com/go-sql-driver/mysql"

	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	Utils "github.com/akgarg0472/urlshortener-auth-service/utils"
)

var (
	logger   = Logger.GetLogger("mysql.go")
	instance *sql.DB
	once     sync.Once
)

func InitDB() {
	once.Do(func() {
		logger.Info("initializing MySQL database")

		db, err := sql.Open("mysql", getDatasource())

		if err != nil {
			logger.Fatal("Error initializing MySQL database: {}", err.Error())
			panic("Error Initializing MySQL Database: " + err.Error())
		}

		pingErr := db.Ping()
		if pingErr != nil {
			logger.Fatal("Ping to DB Failed: {}", pingErr.Error())
			panic("Ping to DB Failed: " + pingErr.Error())
		}

		instance = db
		validateDatabaseSchema()
		initConnectionPool(db)
	})
}

func GetInstance(requestId string, from string) *sql.DB {
	logger.Trace("[{}]: {} getting DB instance", requestId, from)
	return instance
}

func getDatasource() string {
	dbHost := Utils.GetEnvVariable("MYSQL_DB_HOST", "127.0.0.1")
	dbPort := Utils.GetEnvVariable("MYSQL_DB_PORT", "3306")
	dbUserName := Utils.GetEnvVariable("MYSQL_DB_USERNAME", "")
	dbPassword := Utils.GetEnvVariable("MYSQL_DB_PASSWORD", "")
	dbName := Utils.GetEnvVariable("MYSQL_DB_NAME", "")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUserName, dbPassword, dbHost, dbPort, dbName)
}

func initConnectionPool(db *sql.DB) {
	maxIdleConnection, _ := strconv.Atoi(Utils.GetEnvVariable("MYSQL_CONNECTION_POOL_MAX_IDLE_CONNECTION", "10"))
	maxOpenConnection, _ := strconv.Atoi(Utils.GetEnvVariable("MYSQL_CONNECTION_POOL_MAX_OPEN_CONNECTION", "50"))
	db.SetMaxIdleConns(maxIdleConnection)
	db.SetMaxOpenConns(maxOpenConnection)
}

func CloseDB() error {
	if instance != nil {
		err := instance.Close()

		if err != nil {
			logger.Error("Error closing DB: {}", err.Error())
			return err
		}

		logger.Info("DB Instance closed successfully")
	}

	return nil
}

func validateDatabaseSchema() {
	dbName := Utils.GetEnvVariable("MYSQL_DB_NAME", "")
	tableName := Utils.GetEnvVariable("MYSQL_USERS_TABLE_NAME", "")
	schemaValidationQuery := fmt.Sprintf("SELECT count(*) FROM information_schema.tables WHERE table_schema = '%s' AND table_name = '%s';", dbName, tableName)

	var count int = 0
	err := instance.QueryRow(schemaValidationQuery).Scan(&count)

	if err != nil {
		panic(fmt.Sprintf("Error validating DB schema: %s", err.Error()))
	}

	if count != 1 {
		panic(fmt.Sprintf("Expected 1 table with name '%s' but found %d table(s) in '%s' database", tableName, count, dbName))
	}
}
