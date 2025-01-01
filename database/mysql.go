package database

import (
	"fmt"
	"sync"
	"time"

	entity2 "github.com/akgarg0472/urlshortener-auth-service/internal/entity"

	MyLogger "github.com/akgarg0472/urlshortener-auth-service/pkg/logger"
	Utils "github.com/akgarg0472/urlshortener-auth-service/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	ormLogger "gorm.io/gorm/logger"
)

var (
	logger   = MyLogger.GetLogger("mysql.go")
	instance *gorm.DB
	once     sync.Once
)

func InitDB() {
	once.Do(func() {
		logger.Info("initializing MySQL database")

		maxRetryDuration := Utils.GetEnvDurationSeconds("DB_MAX_RETRY_DURATION_SECONDS", 1*time.Minute)
		retryDelay := Utils.GetEnvDurationSeconds("DB_RETRY_DELAY_SECONDS", 5*time.Second)
		var startTime = time.Now()

		var db *gorm.DB
		var err error

		for {
			elapsed := time.Since(startTime)

			if elapsed > maxRetryDuration {
				logger.Fatal("Failed to initialize MySQL database after 1 minute: %s", err.Error())
				panic("Error Initializing MySQL Database: " + err.Error())
			}

			dsn := getDatasource()

			db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
				Logger: ormLogger.Default.LogMode(ormLogger.Silent),
			})

			if err != nil {
				logger.Error("Error initializing MySQL database (elapsed time: %s): %s", elapsed, err.Error())
				time.Sleep(retryDelay)
			} else {
				logger.Info("MySQL database initialized successfully")
				instance = db
				initSchemas()
				return
			}
		}
	})
}

func initSchemas() {
	user := entity2.User{}
	oAuthClient := entity2.OAuthProvider{}

	err := instance.AutoMigrate(&user)

	if err != nil {
		logger.Fatal("Error initializing `{}` database schema: {}", user.TableName(), err.Error())
		panic("Error Initializing DB schema `{}`: " + err.Error())
	}

	err = instance.AutoMigrate(&oAuthClient)

	if err != nil {
		logger.Fatal("Error initializing `{}` database schema: {}", oAuthClient.TableName(), err.Error())
		panic("Error Initializing DB schema: " + err.Error())
	}

	logger.Info("Initialized `{}`, `{}` schema successfully", user.TableName(), oAuthClient.TableName())
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
			logger.Error("Error closing DB: {}", err)
			return err
		}

		logger.Info("DB Instance closed successfully")
	}

	return nil
}
