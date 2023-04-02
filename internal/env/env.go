package env

import (
	"os"
	"strconv"
	"time"
)

func GetString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func MustString(key string) string {
	return getEnvOrPanic(key)
}

func GetBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		if v, err := strconv.ParseBool(value); err == nil {
			return v
		}
	}
	return fallback
}

func MustBool(key string) bool {
	v, err := strconv.ParseBool(getEnvOrPanic(key))
	checkErrAndPanic(err)
	return v
}

func GetInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if v, err := strconv.Atoi(value); err == nil {
			return v
		}
	}
	return fallback
}

func MustInt(key string) int {
	v, err := strconv.Atoi(getEnvOrPanic(key))
	checkErrAndPanic(err)
	return v
}

func GetDuration(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if v, err := time.ParseDuration(value); err == nil {
			return v
		}
	}
	return fallback
}

func MustDuration(key string) time.Duration {
	v, err := time.ParseDuration(getEnvOrPanic(key))
	checkErrAndPanic(err)
	return v
}

// ----------------------------------------------

func getEnvOrPanic(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic("env: " + key + " not found")
	}

	return value
}

func checkErrAndPanic(err error) {
	if err != nil {
		panic(err)
	}
}
