package utils

import (
	"os"
)

func GetEnvVariable(key string, defaultValue string) string {
	value, found := os.LookupEnv(key)

	if !found {
		return defaultValue
	}

	return value
}
