package env

import (
	"os"
	"strconv"
)

func GetString(key, fallback string) string {
	val:= os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
	
}

func GetInt(key string, fallback int) int {
	val:= os.Getenv(key)
	if val == "" {
		return fallback
	}
	
	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return valAsInt
}

func GetBool(key string, fallback bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}

	return boolVal
}