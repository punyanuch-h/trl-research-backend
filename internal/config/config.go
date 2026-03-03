package config

import (
	"log"
	"os"
	"strconv"

	"time"

	"trl-research-backend/internal/utils/send_email"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl            string
	Port             string
	EmailHost        string
	EmailPort        string
	EmailSender      string
	EmailPassword    string
	AntiSpamCooldown string
	JWTExpiry        string
	JWTExpiryTemp    string
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

func (c Config) GetJWTExpiry() time.Duration {
	d, err := time.ParseDuration(c.JWTExpiry)
	if err == nil {
		return d
	}
	// Fallback to integer (now strictly treating as minutes)
	val, err := strconv.Atoi(c.JWTExpiry)
	if err != nil {
		return 480 * time.Minute // Default 8 hours
	}
	return time.Duration(val) * time.Minute
}

func (c Config) GetJWTExpiryTemp() time.Duration {
	d, err := time.ParseDuration(c.JWTExpiryTemp)
	if err == nil {
		return d
	}
	// Fallback to integer (strictly treating as minutes)
	val, err := strconv.Atoi(c.JWTExpiryTemp)
	if err != nil {
		return 10 * time.Minute // Default 10 minutes
	}
	return time.Duration(val) * time.Minute
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
		DBUrl:            dbURL,
		Port:             os.Getenv("PORT"),
		EmailHost:        os.Getenv("EMAIL_HOST"),
		EmailPort:        os.Getenv("EMAIL_PORT"),
		EmailSender:      os.Getenv("EMAIL_USER"),
		EmailPassword:    os.Getenv("EMAIL_PASS"),
		AntiSpamCooldown: getEnv("ANTISPAM_COOLDOWN", "15m"),
		JWTExpiry:        getEnv("JWT_EXPIRY", "8"),
		JWTExpiryTemp:    getEnv("JWT_EXPIRY_TEMP", "10"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
