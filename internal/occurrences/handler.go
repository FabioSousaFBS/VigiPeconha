package occurrences

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(
	service Service,
) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreatePublic(
	c *gin.Context,
) {
	var request CreatePublicOccurrenceRequest

	if err := c.ShouldBindJSON(
		&request,
	); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "dados da ocorrência inválidos",
			},
		)

		return
	}

	occurrence, err :=
		h.service.CreatePublic(
			c.Request.Context(),
			request,
		)

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(
		http.StatusCreated,
		occurrence,
	)
}

func handleError(
	c *gin.Context,
	err error,
) {
	switch {
	case errors.Is(
		err,
		ErrInvalidOccurrenceType,
	),
		errors.Is(
			err,
			ErrAnimalTypeRequired,
		),
		errors.Is(
			err,
			ErrCoordinatesRequired,
		),
		errors.Is(
			err,
			ErrInvalidLatitude,
		),
		errors.Is(
			err,
			ErrInvalidLongitude,
		),
		errors.Is(
			err,
			ErrOccurredAtRequired,
		),
		errors.Is(
			err,
			ErrInvalidState,
		):
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

	default:
		log.Printf(
			"[occurrences] erro interno: %v",
			err,
		)

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "erro interno do servidor",
			},
		)
	}
}
