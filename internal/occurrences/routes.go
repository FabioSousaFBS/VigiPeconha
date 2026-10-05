package occurrences

import (
	"github.com/FabioSousaFBS/vigipeconha-api/internal/contracts"
	"github.com/FabioSousaFBS/vigipeconha-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	handler *Handler,
	jwtSecret string,
	contractService contracts.Service,
) {
	// -----------------------------------------
	// Public routes
	// -----------------------------------------

	public := router.Group("/public")

	public.POST(
		"/occurrences",
		handler.CreatePublic,
	)

	public.POST(
		"/occurrences/:id/photos",
		handler.UploadPhoto,
	)

	// -----------------------------------------
	// Protected routes
	// -----------------------------------------

	protected := router.Group("/occurrences")

	protected.Use(
		middleware.Auth(
			jwtSecret,
		),
		middleware.ActiveContract(
			contractService,
		),
	)

	protected.GET(
		"",
		handler.List,
	)
}
