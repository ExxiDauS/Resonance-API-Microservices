package main

import (
	"auth-service/internal/handler"
	"auth-service/internal/infrastructure/database"
	"auth-service/internal/repository/postgres"
	"auth-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {

	db, err := database.NewPostgresDatabaseClient()
	if err != nil {
		panic(err)
	}

	// migrate tables
	db.AutoMigrate(
		&postgres.User{},
		&postgres.RefreshToken{},
	)

	authService := service.NewAuthService(db)
	authHandler := handler.NewAuthHandler(authService)

	r := gin.Default()

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	r.Run(":8080")
}
