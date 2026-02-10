package track

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetRandomTracks(c *gin.Context) {
	tracks, err := h.service.getRandomFromSpotify(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tracks)
}

func (h *Handler) GetSuggestions(c *gin.Context) {
	source := c.DefaultQuery("source", "spotify") // Default to Spotify
	
	tracks, err := h.service.GetSuggestions(c.Request.Context(), source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tracks)
}