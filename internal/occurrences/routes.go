package occurrences

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	router *gin.Engine,
	handler *Handler,
) {
	public := router.Group("/public")

	public.POST(
		"/occurrences",
		handler.CreatePublic,
	)
}
