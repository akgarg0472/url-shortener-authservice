package logger

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/akgarg0472/urlshortener-auth-service/constants"
	"github.com/akgarg0472/urlshortener-auth-service/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	EnableConsoleLogging bool
	EnableFileLogging    bool
	FileBasePath         string
	LogLevel             zapcore.Level
	EnableStreamLogging  bool
	StreamHost           string
	StreamPort           string
}

var rootLogger *zap.Logger

func init() {
	var err error
	rootLogger, err = createRootLogger()

	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	rootLogger.Info("Logger initialized successfully")
}

// newConfigFromEnv creates a Config instance by reading environment variables.
//
// Environment variables:
//   - SERVICE_NAME: The service name for log tagging.
//   - LOGGING_CONSOLE_ENABLED: Enables/disables console logging.
//   - LOGGING_FILE_ENABLED: Enables/disables file logging.
//   - LOGGING_FILE_BASE_PATH: Specifies the base path for log files.
//   - LOGGING_LEVEL: Defines the logging level (DEBUG, INFO, WARN, ERROR).
//   - LOGGING_STREAM_ENABLED: Enables/disables TCP stream logging.
//   - LOGGING_STREAM_HOST: The hostname for stream logging.
//   - LOGGING_STREAM_PORT: The port for stream logging.
//
// Returns a pointer to the Config struct.
func newConfigFromEnv() *Config {
	parseBool := func(key string) bool {
		v := os.Getenv(key)
		b, err := strconv.ParseBool(v)
		if err != nil {
			return false
		}
		return b
	}

	level := zap.InfoLevel
	if lvlStr := os.Getenv("LOGGING_LEVEL"); lvlStr != "" {
		_ = level.UnmarshalText([]byte(lvlStr))
	}

	return &Config{
		EnableConsoleLogging: parseBool("LOGGING_CONSOLE_ENABLED"),
		EnableFileLogging:    parseBool("LOGGING_FILE_ENABLED"),
		FileBasePath:         os.Getenv("LOGGING_FILE_BASE_PATH"),
		LogLevel:             level,
		EnableStreamLogging:  parseBool("LOGGING_STREAM_ENABLED"),
		StreamHost:           os.Getenv("LOGGING_STREAM_HOST"),
		StreamPort:           os.Getenv("LOGGING_STREAM_PORT"),
	}
}

// createRootLogger initializes and configures the global zap logger based on the settings in the Config struct.
//
// It supports logging to:
//   - Console (stdout) if enabled.
//   - A log file using Lumberjack for log rotation.
//   - A TCP stream if configured.
//
// Returns a zap.Logger instance or an error if the logger cannot be created.
func createRootLogger() (*zap.Logger, error) {
	var cores []zapcore.Core

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "level",
		NameKey:       "logger",
		MessageKey:    "message",
		StacktraceKey: "stackTrace",
		CallerKey:     "caller",
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}

	encoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.000Z"))
	})

	cfg := newConfigFromEnv()

	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
	levelEnabler := zap.NewAtomicLevelAt(cfg.LogLevel)

	if cfg.EnableConsoleLogging {
		consoleCore := zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stdout), levelEnabler)
		cores = append(cores, consoleCore)
	}

	if cfg.EnableFileLogging {
		logFilePath := cfg.FileBasePath + "/" + constants.ServiceName + ".log"
		ljLogger := &lumberjack.Logger{
			Filename:   logFilePath,
			MaxSize:    100,
			MaxBackups: 7,
			MaxAge:     30,
			Compress:   true,
		}
		fileCore := zapcore.NewCore(jsonEncoder, zapcore.AddSync(ljLogger), levelEnabler)
		cores = append(cores, fileCore)
	}

	if cfg.EnableStreamLogging {
		tcpWriter, err := NewTCPAsyncWriter(cfg.StreamHost, cfg.StreamPort)
		if err != nil {
			return nil, err
		}
		streamCore := zapcore.NewCore(jsonEncoder, zapcore.AddSync(tcpWriter), levelEnabler)
		cores = append(cores, streamCore)
	}

	combinedCore := zapcore.NewTee(cores...)

	rootLogger = zap.New(combinedCore, zap.AddCaller(), zap.AddCallerSkip(1)).
		With(
			zap.String(constants.ServiceNameLogKey, constants.ServiceName),
			zap.String(constants.ServiceHostLogKey, utils.GetHostIP()),
		)
	return rootLogger, nil
}

type TCPWriter struct {
	conn net.Conn
}

// NewTCPWriter creates a new TCPWriter instance by establishing a connection to the specified host and port.
//
// Parameters:
//   - host: The hostname or IP address of the log server.
//   - port: The port number of the log server.
//
// Returns a pointer to a TCPWriter instance or an error if the connection fails.
func NewTCPWriter(host string, port string) (*TCPWriter, error) {
	address := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return nil, err
	}
	return &TCPWriter{conn: conn}, nil
}

// Write sends log data to the connected TCP server.
//
// It implements the io.Writer interface, making it compatible with zap logging.
//
// Parameters:
//   - p: The log data to write.
//
// Returns the number of bytes written or an error if the write operation fails.
func (w *TCPWriter) Write(p []byte) (n int, err error) {
	return w.conn.Write(p)
}

// Debug logs a debug-level message.
func Debug(msg string, fields ...zap.Field) {
	rootLogger.Debug(msg, fields...)
}

// Info logs an info-level message.
func Info(msg string, fields ...zap.Field) {
	rootLogger.Info(msg, fields...)
}

// Warn logs a warn-level message.
func Warn(msg string, fields ...zap.Field) {
	rootLogger.Warn(msg, fields...)
}

// Error logs an error-level message.
func Error(msg string, fields ...zap.Field) {
	rootLogger.Error(msg, fields...)
}

// DPanic logs a DPanic-level message.
func DPanic(msg string, fields ...zap.Field) {
	rootLogger.DPanic(msg, fields...)
}

// Panic logs a panic-level message and then panics.
func Panic(msg string, fields ...zap.Field) {
	rootLogger.Panic(msg, fields...)
}

// Fatal logs a fatal-level message and then exits the application.
func Fatal(msg string, fields ...zap.Field) {
	rootLogger.Fatal(msg, fields...)
}

// IsDebugEnabled returns true if the Debug level is enabled.
func IsDebugEnabled() bool {
	return rootLogger.Core().Enabled(zap.DebugLevel)
}

// IsInfoEnabled returns true if the Info level is enabled.
func IsInfoEnabled() bool {
	return rootLogger.Core().Enabled(zap.InfoLevel)
}

// IsWarnEnabled returns true if the Warn level is enabled.
func IsWarnEnabled() bool {
	return rootLogger.Core().Enabled(zap.WarnLevel)
}

// IsErrorEnabled returns true if the Error level is enabled.
func IsErrorEnabled() bool {
	return rootLogger.Core().Enabled(zap.ErrorLevel)
}

// IsDPanicEnabled returns true if the DPanic level is enabled.
func IsDPanicEnabled() bool {
	return rootLogger.Core().Enabled(zap.DPanicLevel)
}

// IsPanicEnabled returns true if the Panic level is enabled.
func IsPanicEnabled() bool {
	return rootLogger.Core().Enabled(zap.PanicLevel)
}

// IsFatalEnabled returns true if the Fatal level is enabled.
func IsFatalEnabled() bool {
	return rootLogger.Core().Enabled(zap.FatalLevel)
}

// Sugar returns a sugared logger for formatted logging.
func Sugar() *zap.SugaredLogger {
	return rootLogger.Sugar()
}

func AddPortToLogger(port int) {
	rootLogger = rootLogger.With(
		zap.Int(constants.ServicePortLogKey, port),
	)
}
