package config

import (
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Port        int
	DatabasePath string
	ProxyPort   int
}

func Load() *Config {
	// Get user's home directory for data storage
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	dataDir := filepath.Join(homeDir, ".quotio")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory %s: %v", dataDir, err)
	}

	return &Config{
		Port:        8080,
		DatabasePath: filepath.Join(dataDir, "quotio.db"),
		ProxyPort:   8081,
	}
}

