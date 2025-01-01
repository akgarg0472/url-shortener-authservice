package utils

import (
	"os"
	"strconv"
	"time"
)

func GetEnvVariable(key string, defaultValue string) string {
	value, found := os.LookupEnv(key)

	if !found {
		return defaultValue
	}

	return value
}

func GetEnvDurationSeconds(envVar string, defaultValue time.Duration) time.Duration {
	envValue := os.Getenv(envVar)

	if envValue == "" {
		return defaultValue
	}

	parsedValue, err := strconv.Atoi(envValue)

	if err != nil {
		return defaultValue
	}

	return time.Duration(parsedValue) * time.Second
}
