package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"music-service/internal/handler"
)

func main() {

	r := gin.Default()
	r.GET("/tracks/random", handler.GetRandomTracks)

	log.Println("Server running on :8080")
	r.Run(":8080")

}
