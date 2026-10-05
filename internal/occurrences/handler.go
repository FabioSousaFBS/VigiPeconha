package occurrences

import (
	"errors"
	"log"
	"net/http"

	"strconv"

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

func (h *Handler) UploadPhoto(c *gin.Context) {
	occurrenceID := c.Param("id")

	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "foto é obrigatória",
			},
		)
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "não foi possível ler a foto",
			},
		)
		return
	}

	defer openedFile.Close()

	photo, err := h.service.UploadPhoto(
		c.Request.Context(),
		occurrenceID,
		file.Size,
		openedFile,
	)

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(
		http.StatusCreated,
		photo,
	)
}

func handleError(
	c *gin.Context,
	err error,
) {
	switch {
	case errors.Is(err, ErrOccurrenceNotFound):
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": err.Error(),
			},
		)

	case errors.Is(err, ErrEmptyPhoto),
		errors.Is(err, ErrPhotoTooLarge),
		errors.Is(err, ErrInvalidPhotoType),
		errors.Is(err, ErrInvalidPhoto),
		errors.Is(err, ErrInvalidOccurrenceType),
		errors.Is(err, ErrInvalidState):

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

func (h *Handler) List(
	c *gin.Context,
) {
	page := parsePositiveInt(
		c.Query("page"),
		1,
	)

	pageSize := parsePositiveInt(
		c.Query("page_size"),
		20,
	)

	filter := ListOccurrencesFilter{
		Page:     page,
		PageSize: pageSize,

		Status: queryStringPointer(
			c,
			"status",
		),

		OccurrenceType: queryStringPointer(
			c,
			"occurrence_type",
		),

		AnimalType: queryStringPointer(
			c,
			"animal_type",
		),

		State: queryStringPointer(
			c,
			"state",
		),

		City: queryStringPointer(
			c,
			"city",
		),
	}

	response, err := h.service.List(
		c.Request.Context(),
		filter,
	)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(
		http.StatusOK,
		response,
	)
}

func parsePositiveInt(
	value string,
	defaultValue int,
) int {
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return defaultValue
	}

	return parsed
}

func queryStringPointer(
	c *gin.Context,
	key string,
) *string {
	value := c.Query(key)

	if value == "" {
		return nil
	}

	return &value
}
