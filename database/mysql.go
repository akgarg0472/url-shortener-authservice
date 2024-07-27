package database

import (
	"fmt"
	"sync"

	entity "github.com/akgarg0472/urlshortener-auth-service/internal/dao/entity"
	Logger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	Utils "github.com/akgarg0472/urlshortener-auth-service/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	logger   = Logger.GetLogger("mysql.go")
	instance *gorm.DB
	once     sync.Once
)

func InitDB() {
	once.Do(func() {
		logger.Info("initializing MySQL database")

		dsn := getDatasource()
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if err != nil {
			logger.Fatal("Error initializing MySQL database: {}", err.Error())
			panic("Error Initializing MySQL Database: " + err.Error())
		}

		instance = db
		initSchema()
	})
}

func initSchema() {
	err := instance.AutoMigrate(&entity.User{})

	if err != nil {
		logger.Fatal("Error initializing database schema: {}", err.Error())
		panic("Error Initializing DB schema: " + err.Error())
	}

	logger.Info("Database schema initialized successfully")
}

func GetInstance(requestId string, from string) *gorm.DB {
	logger.Trace("[{}]: {} getting DB instance", requestId, from)
	return instance
}

func getDatasource() string {
	dbHost := Utils.GetEnvVariable("MYSQL_DB_HOST", "127.0.0.1")
	dbPort := Utils.GetEnvVariable("MYSQL_DB_PORT", "3306")
	dbUserName := Utils.GetEnvVariable("MYSQL_DB_USERNAME", "")
	dbPassword := Utils.GetEnvVariable("MYSQL_DB_PASSWORD", "")
	dbName := Utils.GetEnvVariable("MYSQL_DB_NAME", "")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True", dbUserName, dbPassword, dbHost, dbPort, dbName)
}

func CloseDB() error {
	if instance != nil {
		db, err := instance.DB()

		if err != nil || db.Close() != nil {
			logger.Error("Error closing DB: {}", err.Error())
			return err
		}

		logger.Info("DB Instance closed successfully")
	}

	return nil
}
