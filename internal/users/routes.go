package users

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
	usersGroup := router.Group("/users")

	usersGroup.GET(
		"/me",
		middleware.Auth(jwtSecret),
		handler.Me,
	)

	organizationUsers :=
		router.Group("/organization/users")

	organizationUsers.Use(
		middleware.Auth(jwtSecret),
		middleware.ActiveContract(
			contractService,
		),
		middleware.RequireRole(
			"organization_admin",
		),
	)

	organizationUsers.GET(
		"/license",
		handler.GetLicenseUsage,
	)

	organizationUsers.POST(
		"",
		handler.Create,
	)

	organizationUsers.PATCH(
		"/:id/status",
		handler.UpdateStatus,
	)
}

/*

GET   /users/me

GET   /organization/users/license
POST  /organization/users
PATCH /organization/users/:id/status

*/
