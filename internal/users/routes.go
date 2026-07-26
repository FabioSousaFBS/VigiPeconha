package users

import (
	"github.com/FabioSousaFBS/vigipeconha-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	handler *Handler,
	jwtSecret string,
) {
	usersGroup := router.Group("/users")

	usersGroup.Use(
		middleware.Auth(jwtSecret),
	)

	usersGroup.GET("/me", handler.Me)
}
