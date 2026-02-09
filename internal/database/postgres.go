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
	log.Println("DB_URL: ", dsn)

	if dsn == "" {
		log.Println("❌ DB_URL environment variable is not set")
	}

	var err error
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Use simple protocol to avoid issues with PgBouncer/Neon
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	// ✅ Auto Migrate models
	// err = DB.AutoMigrate(
	// 	&models.Admins{},
	// 	&models.Researchers{},
	// 	&models.Coordinators{},
	// 	&models.Supportments{},
	// 	&models.Cases{},
	// 	&models.Appointments{},
	// 	&models.IntellectualProperties{},
	// 	&models.Assessments{},
	// 	&models.Files{},
	// )
	// if err != nil {
	// 	log.Printf("⚠️ AutoMigrate error: %v", err)
	// }

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
