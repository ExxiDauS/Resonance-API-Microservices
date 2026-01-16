package database

import (
	"log"

	"playlist-service/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializePostgresDB() (*gorm.DB, error) {
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		log.Fatalf("Failed to load database configuration: %v", err)
		return nil, err
	}
	dbURL := dbConfig.DbURL
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to Postgres database: %v", err)
		return nil, err
	}
	return db, nil
}
