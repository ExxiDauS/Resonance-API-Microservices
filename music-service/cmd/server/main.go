package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"music-service/internal/config"
	"music-service/internal/infrastructure/database"
	"music-service/internal/infrastructure/spotify_client"
	"music-service/internal/playlist" // Import your domain package
	"music-service/internal/track"
)

func main() {
	// 1. Config & Infrastructure
	cfg, err := config.LoadPortConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := database.NewPostgresDatabaseClient()
	if err != nil {
		log.Fatal(err)
	}

	// 1. Migrations (Add Playlist)
	db.AutoMigrate(
		&track.Track{}, &track.Artist{}, &track.Genre{}, // Existing
		&playlist.Playlist{}, // <--- NEW
	)

	// 2. Music Module
	spotifyClient, _ := spotify_client.NewSpotifyClient()
	musicRepo := track.NewRepository(db)
	musicService := track.NewService(musicRepo, spotifyClient)
	musicHandler := track.NewHandler(musicService)

	// 3. Playlist Module
	playlistRepo := playlist.NewRepository(db)
	
	// Adapter: MusicService satisfies Playlist's "MusicProvider" interface
	// because MusicService has GetTrack(id) (Wrapped around GetTrackByID)
	playlistService := playlist.NewService(playlistRepo, musicService) 
	// Note: You might need to make sure MusicService or MusicRepo exposes a method 
	// signature that matches exactly what PlaylistService expects.
	
	playlistHandler := playlist.NewHandler(playlistService)

	// 4. Routes
	r := gin.Default()

	// Music Routes
	// ?source=db OR ?source=spotify
	r.GET("/tracks/random", musicHandler.GetRandomTracks)
	r.GET("/tracks/suggestions", musicHandler.GetSuggestions) 

	// Playlist Routes
	r.POST("/playlists", playlistHandler.CreatePlaylist)
	r.GET("/playlists/:id", playlistHandler.GetPlaylist)
	r.POST("/playlists/:id/tracks", playlistHandler.AddTrack)

	r.Run(":" + cfg.Port)
}


