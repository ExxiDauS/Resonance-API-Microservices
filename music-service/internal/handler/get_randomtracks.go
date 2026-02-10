package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"music-service/internal/service"
)

func GetRandomTracks(c *gin.Context) {
	tracks, err := service.GetRandomTracks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tracks)
}
