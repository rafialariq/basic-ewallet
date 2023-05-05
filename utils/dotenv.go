package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func DotEnv(key string, envFilePath string) string {
	if err := godotenv.Load(envFilePath); err != nil {
		log.Fatalln("Error loading .env file")
	}

	return os.Getenv(key)
}
