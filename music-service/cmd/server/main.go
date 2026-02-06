package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"music-service/internal/config"
	"music-service/internal/handler"
	"music-service/internal/service"

	"music-service/internal/infrastructure/database"
	"music-service/internal/repository/postgres"
)

func main() {

	config, err := config.LoadPortConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := database.NewPostgresDatabaseClient()
	err = db.AutoMigrate(
		&postgres.UserGenreScore{},
	)
	if err != nil {
		log.Fatalf("Migrate error: %v", err)
	}

	r := gin.Default()
	r.GET("/tracks/random", handler.GetRandomTracks)

	repo := postgres.NewPreferenceRepository(db)
	svc := service.NewPreferenceService(repo)
	h := handler.NewPreferenceHandler(svc)

	r.POST("/user/preferences", h.Save)

	log.Printf("Server running on :%s", config.Port)
	r.Run(":" + config.Port)

}
