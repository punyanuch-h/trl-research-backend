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

	// init Firestore
	database.InitFirebase("localServiceAccountKey.json")
	defer database.CloseFirebase()

	// Initialize GCSClient
	bucket := os.Getenv("GCS_BUCKET_NAME")
	saEmail := os.Getenv("SA_EMAIL")

	gcsClient, err := storage.NewGCSClient(bucket, saEmail)
	if err != nil {
		log.Fatalf("Failed to initialize GCS client: %v", err)
	}

	// pass gcsClient here
	r := router.SetupRouter(gcsClient)

	// Run server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	fmt.Println("🚀 Server running on", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
