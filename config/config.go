package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	} else {
		log.Println(".env loaded successfully")
	}

	log.Println("DB User:", os.Getenv("DB_USER"))
}
