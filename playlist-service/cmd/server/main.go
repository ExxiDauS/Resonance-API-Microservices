package main

import (
	"log"
	"playlist-service/internal/infrastructure/database"
)

func main() {
	log.Println("Starting Playlist Service...")

	db, err := database.InitializePostgresDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	log.Println("Database connected successfully:", db)
}
