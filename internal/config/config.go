package config

import (
	"log"
	"os"

	"time"

	"trl-research-backend/internal/utils/send_email"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl         string
	Port          string
	EmailHost     string
	EmailPort     string
	EmailSender   string
	EmailPassword    string
	AntiSpamCooldown string
}

func (c Config) GetSMTPConfig() send_email.SMTPConfig {
	return send_email.SMTPConfig{
		Host:     c.EmailHost,
		Port:     c.EmailPort,
		Username: c.EmailSender,
		Password: c.EmailPassword,
		From:     c.EmailSender,
	}
}

func (c Config) GetAntiSpamCooldown() time.Duration {
	d, err := time.ParseDuration(c.AntiSpamCooldown)
	if err != nil {
		return 15 * time.Minute // Default fall back
	}
	return d
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
		AntiSpamCooldown: getEnv("ANTISPAM_COOLDOWN", "15m"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
