package users

import (
	"net/http"

	"github.com/FabioSousaFBS/vigipeconha-api/internal/shared"
	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Me(c *gin.Context) {
	userID, exists := c.Get(
		shared.ContextUserIDKey,
	)

	if !exists {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "usuário não autenticado",
			},
		)
		return
	}

	email, _ := c.Get(
		shared.ContextUserEmailKey,
	)

	role, _ := c.Get(
		shared.ContextUserRoleKey,
	)

	c.JSON(
		http.StatusOK,
		gin.H{
			"id":    userID,
			"email": email,
			"role":  role,
		},
	)
}
