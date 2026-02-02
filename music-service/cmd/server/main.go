package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"music-service/internal/config"
	"music-service/internal/handler"
)

func main() {

	config, err := config.LoadPortConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	r := gin.Default()
	r.GET("/tracks/random", handler.GetRandomTracks)

	log.Printf("Server running on :%s", config.Port)
	r.Run(":" + config.Port)

}
