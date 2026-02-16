package main

import (
	"fmt"
	"log"
	"trl-research-backend/internal/config"
	"trl-research-backend/internal/database"
	"trl-research-backend/internal/models"
)

func main() {
	// Load environment variables
	config.LoadConfig()

	// Initialize Postgres
	database.InitPostgres()
	defer database.ClosePostgres()

	fmt.Println("⏳ Running Migrations...")

	// AutoMigrate all models
	err := database.DB.AutoMigrate(
		&models.Admins{},
		&models.Researchers{},
		&models.Coordinators{},
		&models.Cases{},
		&models.Supportments{},
		&models.Appointments{},
		&models.IntellectualProperties{},
		&models.Assessments{},
		&models.ChatLogs{},
	)

	if err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	fmt.Println("✅ Migrations completed successfully!")
}
