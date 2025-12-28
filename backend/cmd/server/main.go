package main

import (
	"log"
	"quotio-electron-go/backend/internal/api"
	"quotio-electron-go/backend/internal/config"
	"quotio-electron-go/backend/internal/storage"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := storage.Initialize(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize API server
	server := api.NewServer(db, cfg)
	
	// Start server
	log.Printf("Starting server on port %d", cfg.Port)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

