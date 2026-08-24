package middleware

import (
	"net/http"

	"github.com/FabioSousaFBS/vigipeconha-api/internal/contracts"
	"github.com/FabioSousaFBS/vigipeconha-api/internal/shared"
	"github.com/gin-gonic/gin"
)

func ActiveContract(
	service contracts.Service,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get(
			shared.ContextOrganizationIDKey,
		)

		if !exists {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{
					"error": "usuário não possui organização",
				},
			)
			return
		}

		organizationID, ok := value.(string)

		if !ok || organizationID == "" {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{
					"error": "organização inválida",
				},
			)
			return
		}

		err := service.ValidateOrganizationAccess(
			c.Request.Context(),
			organizationID,
		)

		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{
					"error": err.Error(),
				},
			)
			return
		}

		c.Next()
	}
}
