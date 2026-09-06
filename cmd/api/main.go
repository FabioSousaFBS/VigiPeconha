package main

import (
	"context"
	"log"
	"net/http"

	"github.com/FabioSousaFBS/vigipeconha-api/internal/auth"
	"github.com/FabioSousaFBS/vigipeconha-api/internal/config"
	"github.com/FabioSousaFBS/vigipeconha-api/internal/contracts"
	"github.com/FabioSousaFBS/vigipeconha-api/internal/database"
	"github.com/FabioSousaFBS/vigipeconha-api/internal/middleware"
	"github.com/FabioSousaFBS/vigipeconha-api/internal/organizations"
	"github.com/FabioSousaFBS/vigipeconha-api/internal/users"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	db := database.Connect(
		ctx,
		cfg,
	)

	defer db.Close()

	router := gin.Default()

	// -----------------------------------------
	// Organizations
	// -----------------------------------------

	organizationRepository :=
		organizations.NewPostgresRepository(
			db,
		)

	// -----------------------------------------
	// Contracts
	// -----------------------------------------

	contractRepository :=
		contracts.NewPostgresRepository(
			db,
		)

	contractService :=
		contracts.NewService(
			contractRepository,
			organizationRepository,
		)

	// -----------------------------------------
	// Authentication
	// -----------------------------------------

	authRepository :=
		auth.NewPostgresRepository(
			db,
		)

	authService :=
		auth.NewService(
			authRepository,
			cfg.JWTSecret,
		)

	authHandler :=
		auth.NewHandler(
			authService,
		)

	auth.RegisterRoutes(
		router,
		authHandler,
	)

	// -----------------------------------------
	// Users
	// -----------------------------------------

	usersRepository :=
		users.NewPostgresRepository(
			db,
		)

	usersService :=
		users.NewService(
			usersRepository,
		)

	usersHandler :=
		users.NewHandler(
			usersService,
		)

	users.RegisterRoutes(
		router,
		usersHandler,
		cfg.JWTSecret,
		contractService,
	)

	// -----------------------------------------
	// Health
	// -----------------------------------------

	router.GET(
		"/health",
		func(c *gin.Context) {
			err := db.Ping(
				c.Request.Context(),
			)

			if err != nil {
				c.JSON(
					http.StatusInternalServerError,
					gin.H{
						"status":   "error",
						"database": "disconnected",
					},
				)

				return
			}

			c.JSON(
				http.StatusOK,
				gin.H{
					"status":   "ok",
					"database": "connected",
				},
			)
		},
	)

	// -----------------------------------------
	// Protected Test
	// -----------------------------------------

	protected :=
		router.Group(
			"/protected",
		)

	protected.Use(
		middleware.Auth(
			cfg.JWTSecret,
		),
		middleware.ActiveContract(
			contractService,
		),
	)

	protected.GET(
		"/test",
		func(c *gin.Context) {
			c.JSON(
				http.StatusOK,
				gin.H{
					"status": "access_granted",
				},
			)
		},
	)

	// -----------------------------------------
	// Server
	// -----------------------------------------

	if err := router.Run(
		":" + cfg.AppPort,
	); err != nil {
		log.Fatalf(
			"erro ao iniciar servidor: %v",
			err,
		)
	}
}
