package utils

import (
	"os"
	"strconv"
)

func GetIntEnv(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
