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
		&models.AdminInfo{},
		&models.ResearcherInfo{},
		&models.CoordinatorInfo{},
		&models.CaseInfo{},
		&models.Supporter{},
		&models.Appointment{},
		&models.IntellectualProperty{},
		&models.AssessmentTrl{},
		&models.FileMetadata{},
	)

	if err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	fmt.Println("✅ Migrations completed successfully!")
}
