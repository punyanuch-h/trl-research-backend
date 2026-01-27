package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitPostgres() {
	dsn := os.Getenv("DEV_DB_URL") // local
	if dsn == "" {
		dsn = os.Getenv("DB_URL") // production
	}

	if dsn == "" {
		log.Println("❌ DB_URL environment variable is not set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Println("❌ Failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Println("❌ Failed to get database instance: %v", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("✅ Connected to Neon Postgres database")
}

func ClosePostgres() {
	if DB == nil {
		return
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("⚠️ Error getting database instance during close: %v", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Printf("⚠️ Error closing database connection: %v", err)
	}
}
