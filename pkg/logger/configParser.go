package logger

import (
	"strings"

	"github.com/akgarg0472/urlshortener-auth-service/utils"
)

func ReadConfig(path string) Config {
	config := Config{}

	config.Level = utils.GetEnvVariable("LOGGER_LEVEL", "INFO")
	processLogLevel(&config)

	loggerType := utils.GetEnvVariable("LOGGER_TYPE", "console")
	processLoggerType(strings.TrimSpace(loggerType), &config)

	config.Enabled = utils.GetEnvVariable("LOGGER_ENABLED", "true") == "true"

	config.LogFilePath = utils.GetEnvVariable("LOGGER_LOG_FILE_PATH", "/tmp/logs.log")

	return config
}

func processLogLevel(config *Config) {
	switch config.Level {
	case "fatal":
	case "FATAL":
		config.IsFatalEnabled = true

	case "error":
	case "ERROR":
		config.IsErrorEnabled = true
		config.IsFatalEnabled = true

	case "info":
	case "INFO":
		config.IsInfoEnabled = true
		config.IsErrorEnabled = true
		config.IsFatalEnabled = true

	case "debug":
	case "DEBUG":
		config.IsDebugEnabled = true
		config.IsInfoEnabled = true
		config.IsErrorEnabled = true
		config.IsFatalEnabled = true

	case "trace":
	case "TRACE":
		config.IsTraceEnabled = true
		config.IsDebugEnabled = true
		config.IsInfoEnabled = true
		config.IsErrorEnabled = true
		config.IsFatalEnabled = true
	}
}

func processLoggerType(value string, config *Config) {
	switch value {
	case "console":
		config.LogToConsole = true
	case "file":
		config.LogToFile = true
	}
}
