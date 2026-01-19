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
	dsn := os.Getenv("DATABASE_URL")
	log.Fatal("dsn database_url", dsn)
	if dsn == "" {
		dsn = os.Getenv("DB_URL")
		log.Fatal("dsn db_url", dsn)
	}

	if dsn == "" {
		log.Fatal("❌ ")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get database instance: %v", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("✅ Connected to Neon Postgres database")
}

func ClosePostgres() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("⚠️ Error getting database instance during close: %v", err)
		return
	}
	sqlDB.Close()
}
