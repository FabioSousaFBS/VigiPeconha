package healthfacilities

import (
	"github.com/FabioSousaFBS/vigipeconha-api/internal/contracts"
	"github.com/FabioSousaFBS/vigipeconha-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, h *Handler, jwtSecret string, contractService contracts.Service) {
	public := router.Group("/public/health-facilities")
	public.GET("", h.List)
	public.GET("/nearby", h.List)
	public.GET("/:id", h.GetByID)
	protected := router.Group("/health-facilities")
	protected.Use(middleware.Auth(jwtSecret), middleware.ActiveContract(contractService), middleware.RequireRole("organization_admin", "system_admin"))
	protected.POST("", h.Create)
	protected.PATCH("/:id", h.Update)
}
