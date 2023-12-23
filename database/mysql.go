package database

import (
	"database/sql"
	"fmt"
	"os"
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
		fmt.Println("initializing DB")

		db, err := sql.Open("mysql", getDatasource())

		if err != nil {
			logger.Fatal("Error initializing MySQL database: {}", err.Error())
			panic("Error Initializing MySQL Database: " + err.Error())
		}

		initConnectionPool(db)

		pingErr := db.Ping()

		if pingErr != nil {
			logger.Fatal("Ping to DB Failed: {}", pingErr.Error())
			panic("Ping to DB Failed: " + pingErr.Error())
		}

		instance = db

		initSchemaErr := initDatabaseSchema()

		if initSchemaErr != nil {
			logger.Fatal("Schema initialization failed: {}", initSchemaErr.Error())
			panic("Schema initialization failed: " + initSchemaErr.Error())
		}
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

func initDatabaseSchema() error {
	createSQLQueries, err := os.ReadFile("database/queries/create_tables.sql")

	if err != nil {
		return fmt.Errorf("Error reading sql file: " + err.Error())
	}

	result, err := instance.Exec(string(createSQLQueries))

	if err != nil {
		return fmt.Errorf("Error creating tables: " + err.Error())
	}

	logger.Info("Tables initialized successfully with result: {}", result)

	return nil
}
