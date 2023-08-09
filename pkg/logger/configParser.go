package logger

import (
	"bufio"
	"os"
	"strings"
)

func ReadConfig(path string) LoggerConfig {
	if path == "" {
		path = "logger.conf"
	}

	file, err := os.Open(path)

	if err != nil {
		panic("Error reading logger config file")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	config := LoggerConfig{}

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || line[0] == '#' {
			continue
		}

		keyValuePair := strings.Split(line, "=")

		if len(keyValuePair) != 2 {
			panic("Invalid logger config property %s" + line)
		}

		key := strings.TrimSpace(keyValuePair[0])
		value := strings.TrimSpace(keyValuePair[1])

		switch key {
		case "logger.level":
			config.Level = value
			setLogLevel(&config)

		case "logger.type":
			loggerTypes := strings.Split(value, ",")
			for _, loggerType := range loggerTypes {
				handleType(strings.TrimSpace(loggerType), &config)
			}

		case "logger.enabled":
			config.Enabled = value == "true"

		case "logger.filepath":
			config.LogFilePath = value
		}
	}

	return config
}

func setLogLevel(config *LoggerConfig) {
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

func handleType(value string, config *LoggerConfig) {
	switch value {
	case "console":
		config.LogToConsole = true
	case "file":
		config.LogToFile = true
	}
}
