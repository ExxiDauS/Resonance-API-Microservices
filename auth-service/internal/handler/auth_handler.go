package handler

import (
	"net/http"

	"auth-service/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(a *service.AuthService) *AuthHandler {
	return &AuthHandler{a}
}

func (h *AuthHandler) Register(c *gin.Context) {

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	if c.BindJSON(&req) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	err := h.auth.Register(req.Email, req.Password, req.Name)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "registered"})
}

func (h *AuthHandler) Login(c *gin.Context) {

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if c.BindJSON(&req) != nil {
		c.JSON(400, gin.H{"error": "bad request"})
		return
	}

	user, token, err := h.auth.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"access_token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.DisplayName,
		},
	})
}

func (h *AuthHandler) Me(c *gin.Context) {

	userID, _ := c.Get("user_id")

	c.JSON(200, gin.H{
		"user_id": userID,
	})
}
