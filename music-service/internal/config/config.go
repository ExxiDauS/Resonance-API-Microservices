package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	DbURL string
}

type SpotifyConfig struct {
	SpotifyClientID     string
	SpotifyClientSecret string
}

func LoadDatabaseConfig() (*DatabaseConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &DatabaseConfig{
		DbURL: os.Getenv("DATABASE_URL"),
	}

	return config, nil
}

func LoadSpotifyConfig() (*SpotifyConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &SpotifyConfig{
		SpotifyClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		SpotifyClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
	}

	return config, nil
}
