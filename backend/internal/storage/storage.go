package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"path/filepath"
)

var DB *gorm.DB

// Initialize sets up the database connection, encryption, and runs migrations
func Initialize(dbPath string) (*gorm.DB, error) {
	// Initialize encryption first (extract directory from dbPath)
	dataDir := filepath.Dir(dbPath)
	if err := InitEncryption(dataDir); err != nil {
		return nil, err
	}

	var err error

	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}

	// Auto-migrate tables
	err = DB.AutoMigrate(
		&Account{},
		&QuotaHistory{},
		&ProxyConfig{},
		&AgentConfig{},
		&ProviderHealth{},
	)

	if err != nil {
		return nil, err
	}

	// Initialize default proxy config if not exists
	var proxyConfig ProxyConfig
	if err := DB.First(&proxyConfig).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			defaultConfig := ProxyConfig{
				Port:            8081,
				RoutingStrategy: "round_robin",
				AutoStart:       false,
			}
			DB.Create(&defaultConfig)
		}
	}

	log.Println("Database initialized successfully")
	return DB, nil
}
