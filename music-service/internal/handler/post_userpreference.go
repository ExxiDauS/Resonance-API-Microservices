package handler

import (
	"net/http"

	"music-service/internal/service"

	"github.com/gin-gonic/gin"
)

type PreferenceHandler struct {
	service *service.PreferenceService
}

func NewPreferenceHandler(s *service.PreferenceService) *PreferenceHandler {
	return &PreferenceHandler{service: s}
}

type PreferenceRequest struct {
	Genres map[string]int `json:"genres"`
}

func (h *PreferenceHandler) Save(c *gin.Context) {

	userID := c.GetString("user_id")
	var req PreferenceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	err := h.service.Save(userID, req.Genres)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "saved"})

}
