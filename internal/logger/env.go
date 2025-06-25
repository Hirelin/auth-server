package logger

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	env := godotenv.Load()

	if env != nil {
		logger.Fatal("Error loading .env file")
	}
}

func GetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		logger.Fatalf("Environment variable %s not set", key)
	}
	return value
}
