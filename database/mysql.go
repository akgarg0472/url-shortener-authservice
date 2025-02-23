package database

import (
	"fmt"
	"sync"
	"time"

	"github.com/akgarg0472/urlshortener-auth-service/constants"
	"github.com/akgarg0472/urlshortener-auth-service/internal/entity"
	"github.com/akgarg0472/urlshortener-auth-service/internal/logger"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	instance *gorm.DB
	once     sync.Once
)

func InitDB() {
	once.Do(func() {
		logger.Info("initializing MySQL database")

		maxRetryDuration := utils.GetEnvDurationSeconds("DB_MAX_RETRY_DURATION_SECONDS", 1*time.Minute)
		retryDelay := utils.GetEnvDurationSeconds("DB_RETRY_DELAY_SECONDS", 5*time.Second)
		var startTime = time.Now()

		var db *gorm.DB
		var err error

		for {
			elapsed := time.Since(startTime)

			if elapsed > maxRetryDuration {
				if logger.IsFatalEnabled() {
					logger.Fatal("Failed to initialize MySQL database after 1 minute", zap.Error(err))
				}
				panic(fmt.Sprintf("Error Initializing MySQL Database: %v", err))
			}

			dsn := getDatasource()

			db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
				Logger: gormLogger.Default.LogMode(gormLogger.Silent),
			})

			if err != nil {
				if logger.IsErrorEnabled() {
					logger.Error("Error initializing MySQL database",
						zap.Duration("elapsed_time", elapsed),
						zap.Error(err),
					)
				}
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
	user := entity.User{}
	oAuthClient := entity.OAuthProvider{}

	err := instance.AutoMigrate(&user)

	if err != nil {
		if logger.IsFatalEnabled() {
			logger.Fatal("Error initializing database schema",
				zap.String("table", user.TableName()),
				zap.Error(err),
			)
		}
		panic(fmt.Sprintf("Error initializing DB schema `%s`: %v", user.TableName(), err))
	}

	err = instance.AutoMigrate(&oAuthClient)

	if err != nil {
		if logger.IsFatalEnabled() {
			logger.Fatal("Error initializing database schema", zap.String("schema", oAuthClient.TableName()), zap.Error(err))
		}
		panic(fmt.Sprintf("Error initializing DB schema `%s`: %v", oAuthClient.TableName(), err))
	}

	logger.Info("Initialized database schemas successfully",
		zap.String("user_schema", user.TableName()),
		zap.String("oauth_client_schema", oAuthClient.TableName()),
	)
}

func GetInstance(requestId string, from string) *gorm.DB {
	if logger.IsDebugEnabled() {
		logger.Debug("Getting DB instance",
			zap.String(constants.RequestIdLogKey, requestId),
			zap.String("from", from),
		)
	}
	return instance
}

func getDatasource() string {
	dbHost := utils.GetEnvVariable("MYSQL_DB_HOST", "127.0.0.1")
	dbPort := utils.GetEnvVariable("MYSQL_DB_PORT", "3306")
	dbUserName := utils.GetEnvVariable("MYSQL_DB_USERNAME", "")
	dbPassword := utils.GetEnvVariable("MYSQL_DB_PASSWORD", "")
	dbName := utils.GetEnvVariable("MYSQL_DB_NAME", "")
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True", dbUserName, dbPassword, dbHost, dbPort, dbName)
}

func CloseDB() error {
	if instance != nil {
		db, err := instance.DB()

		if err != nil || db.Close() != nil {
			if logger.IsErrorEnabled() {
				logger.Error("Error closing DB", zap.Error(err))
			}
			return err
		}

		logger.Info("DB Instance closed successfully")
	}

	return nil
}
