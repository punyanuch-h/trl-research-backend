package main

import (
	"fmt"
	"log"
	"os"

	"trl-research-backend/internal/config"
	"trl-research-backend/internal/database"
	"trl-research-backend/internal/router"
)

func main() {
	// Load environment variables from .env file
	config.LoadConfig()

	// on cloud
	// database.InitFirebase("serviceAccountKey.json")
	// local
	database.InitFirebase("trl-research-service-account.json")
	defer database.CloseFirebase()

	// Setup router
	r := router.SetupRouter()

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
