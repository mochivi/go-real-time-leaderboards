package utils

import (
	"log"
	"os"
	"strconv"
)

func GetEnvString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("Env var %s not set, using fallback: %s", key, fallback)
		return fallback
	}

	return val
}

func GetEnvInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("Env var %s not set, using fallback: %d", key, fallback)
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Env var %s not an integer, using fallback: %d", key, fallback)
		return fallback
	}

	return valAsInt
}

func GetEnvBool(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("Env var %s not set, using fallback: %t", key, fallback)
		return fallback
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		log.Printf("Env var %s not a boolean, using fallback: %t", key, fallback)
		return fallback
	}

	return boolVal
}