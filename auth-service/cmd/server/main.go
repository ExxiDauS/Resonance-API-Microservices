package main

import (
	"auth-service/internal/handler"
	"auth-service/internal/infrastructure/database"
	"auth-service/internal/middleware"
	"auth-service/internal/repository/postgres"
	"auth-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {

	db, err := database.NewPostgresDatabaseClient()
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&postgres.User{},
		&postgres.RefreshToken{},
	)

	authService := service.NewAuthService(db)
	authHandler := handler.NewAuthHandler(authService)

	r := gin.Default()

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	protected := r.Group("/")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/me", authHandler.Me)
	}
	r.Run(":9090")
}
