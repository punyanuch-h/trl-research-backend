package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl         string
	Port          string
	EmailHost     string
	EmailPort     string
	EmailSender   string
	EmailPassword string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No .env file found, using system environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DEV_DB_URL") // local
	}

	if dbURL == "" {
		log.Fatal("❌ DB_URL is not set")
	}

	return Config{
		DBUrl:         dbURL,
		Port:          os.Getenv("PORT"),
		EmailHost:     os.Getenv("EMAIL_HOST"),
		EmailPort:     os.Getenv("EMAIL_PORT"),
		EmailSender:   os.Getenv("EMAIL_SENDER"),
		EmailPassword: os.Getenv("EMAIL_PASSWORD"),
	}
}
