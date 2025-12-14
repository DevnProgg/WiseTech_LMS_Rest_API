package main

import (
	"log"

	"wisetech-lms-api/internal/config"
	"wisetech-lms-api/internal/database"
	"wisetech-lms-api/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize database schema
	if err := database.InitializeSchema(db); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	// Create a new server
	srv := server.New(db, cfg)

	// Start the server
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
