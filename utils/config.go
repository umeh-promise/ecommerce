package utils

import (
	"os"
	"strconv"
)

type Config struct {
	Addr        string
	MaxOpenConn int
	MaxIdleConn int
	MaxIdleTime string
}

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}

	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	valueAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return valueAsInt
}
