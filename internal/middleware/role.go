package middleware

import (
	"net/http"

	"github.com/FabioSousaFBS/vigipeconha-api/internal/shared"
	"github.com/gin-gonic/gin"
)

func RequireRole(
	allowedRoles ...string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get(
			shared.ContextUserRoleKey,
		)

		if !exists {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{
					"error": "perfil do usuário não encontrado",
				},
			)
			return
		}

		role, ok := value.(string)

		if !ok || role == "" {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{
					"error": "perfil do usuário inválido",
				},
			)
			return
		}

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{
				"error": "usuário não possui permissão para esta operação",
			},
		)
	}
}
