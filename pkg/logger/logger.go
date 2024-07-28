package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var loggerConfig Config
var replacer = strings.NewReplacer("{}", "%v")

type Logger struct {
	config Config
	file   string
}

func init() {
	loggerConfig = ReadConfig("logger.conf")

	if loggerConfig.Enabled {
		initLogger(&loggerConfig)
	}
}

func initLogger(config *Config) {
	log.SetFlags(log.Ldate | log.Ltime)

	if config.LogToFile {
		file, err := os.OpenFile(config.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

		if err != nil {
			panic("Error creating/opening log file:" + err.Error())
		}

		log.SetOutput(file)
	}
}

func GetLogger(file string) *Logger {
	return &Logger{
		config: loggerConfig,
		file:   file,
	}
}

func (l *Logger) Info(message string, args ...interface{}) {
	if l.config.IsInfoEnabled {
		if len(args) > 0 {
			message = fmt.Sprintf(replacer.Replace(message), args...)
		}
		doLog(l.file, "INFO", message)
	}
}

func (l *Logger) Error(message string, args ...interface{}) {
	if l.config.IsErrorEnabled {
		if len(args) > 0 {
			message = fmt.Sprintf(replacer.Replace(message), args...)
		}
		doLog(l.file, "ERROR", message)
	}
}

func (l *Logger) Fatal(message string, args ...interface{}) {
	if l.config.IsFatalEnabled {
		if len(args) > 0 {
			message = fmt.Sprintf(replacer.Replace(message), args...)
		}
		doLog(l.file, "FATAL", message)
	}
}

func (l *Logger) Debug(message string, args ...interface{}) {
	if l.config.IsDebugEnabled {
		if len(args) > 0 {
			message = fmt.Sprintf(replacer.Replace(message), args...)
		}
		doLog(l.file, "DEBUG", message)
	}
}

func (l *Logger) Trace(message string, args ...interface{}) {
	if l.config.IsTraceEnabled {
		if len(args) > 0 {
			message = fmt.Sprintf(replacer.Replace(message), args...)
		}
		doLog(l.file, "TRACE", message)
	}
}

func doLog(file string, level string, message string) {
	if loggerConfig.LogToFile {
		log.Println("[" + level + "] " + file + " - " + message)
	}

	if loggerConfig.LogToConsole {
		currentTime := getCurrentTimestamp()
		println(currentTime + " [" + level + "] " + file + " - " + message)
	}
}

func getCurrentTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
