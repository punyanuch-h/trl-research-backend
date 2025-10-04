package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl string
	Port  string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No .env file found, using system environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		// Support Neon default env var name
		dbURL = os.Getenv("DATABASE_URL")
	}

	return Config{
		DBUrl: dbURL,
		Port:  os.Getenv("PORT"),
	}
}
