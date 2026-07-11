package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"vessel.dev/vessel/internal/cloud/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load")
	}

	port := os.Getenv("CLOUD_PORT")
	if port == "" {
		port = "8081" // Different from OSS vessel default (8080)
	}

	db := server.InitDatabase()

	app := server.NewServer(db)
	log.Printf("Starting Vessel Cloud API on port %s", port)
	if err := app.Start(":" + port); err != nil {
		log.Fatalf("Failed to start cloud server: %v", err)
	}
}
