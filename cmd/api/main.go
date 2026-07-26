package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/FabioSousaFBS/vigipeconha-api/internal/auth"
	"github.com/FabioSousaFBS/vigipeconha-api/internal/config"
	"github.com/FabioSousaFBS/vigipeconha-api/internal/database"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	db := database.Connect(ctx, cfg)
	defer db.Close()

	router := gin.Default()

	authRepository := auth.NewPostgresRepository(db)
	authService := auth.NewService(
		authRepository,
		cfg.JWTSecret,
	)
	authHandler := auth.NewHandler(authService)

	auth.RegisterRoutes(router, authHandler)

	router.GET("/health", func(c *gin.Context) {
		err := db.Ping(ctx)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":   "error",
				"database": "disconnected",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"database": "connected",
		})
	})

	router.Run(fmt.Sprintf(":%s", cfg.AppPort))
}
