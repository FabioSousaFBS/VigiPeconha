package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	router *gin.Engine,
	handler *Handler,
) {
	authGroup := router.Group("/auth")

	authGroup.POST("/register", handler.Register)
	authGroup.POST("/login", handler.Login)
}
