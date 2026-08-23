package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}
}

func GetEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		log.Printf(
			"invalid value for %s=%q, using default %d",
			key,
			value,
			defaultValue,
		)
		return defaultValue
	}

	return result
}

func GetEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	result, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf(
			"invalid value for %s=%q, using default %t",
			key,
			value,
			defaultValue,
		)
		return defaultValue
	}

	return result
}
