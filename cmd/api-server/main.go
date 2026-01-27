package main

import (
	"fmt"
	"log"
	"os"

	"trl-research-backend/internal/config"
	"trl-research-backend/internal/database"
	"trl-research-backend/internal/router"
	"trl-research-backend/internal/storage"
)

func main() {
	// Load environment variables (.env)
	config.LoadConfig()

	// Initialize Postgres (Neon)
	database.InitPostgres()
	defer database.ClosePostgres()

	// Initialize GCSClient
	bucket := os.Getenv("GCS_BUCKET_NAME")
	saEmail := os.Getenv("SA_EMAIL")
	log.Println("GCS_BUCKET_NAME:", bucket)
	log.Println("SA_EMAIL:", saEmail)

	var gcsClient *storage.GCSClient
	if bucket != "" {
		client, err := storage.NewGCSClient(bucket, saEmail)
		if err != nil {
			log.Println("⚠️ GCS init failed, continue without GCS:", err)
		} else {
			gcsClient = client
		}
	} else {
		log.Println("⚠️ GCS_BUCKET_NAME not set, skip GCS")
	}

	// pass gcsClient here
	r := router.SetupRouter(gcsClient)

	// Run server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf("0.0.0.0:%s", port)

	fmt.Println("🚀 Server running on", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
