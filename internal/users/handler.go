package users

import (
	"errors"
	"log"
	"net/http"

	"github.com/FabioSousaFBS/vigipeconha-api/internal/shared"
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

func (h *Handler) Me(
	c *gin.Context,
) {
	userID, _ := c.Get(
		shared.ContextUserIDKey,
	)

	email, _ := c.Get(
		shared.ContextUserEmailKey,
	)

	role, _ := c.Get(
		shared.ContextUserRoleKey,
	)

	organizationID, _ := c.Get(
		shared.ContextOrganizationIDKey,
	)

	response := gin.H{
		"id":    userID,
		"email": email,
		"role":  role,
	}

	if organizationID != nil {
		response["organization_id"] =
			organizationID
	}

	c.JSON(
		http.StatusOK,
		response,
	)
}

func (h *Handler) GetLicenseUsage(
	c *gin.Context,
) {
	organizationID, ok :=
		getOrganizationID(c)

	if !ok {
		return
	}

	usage, err := h.service.GetLicenseUsage(
		c.Request.Context(),
		organizationID,
	)

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(
		http.StatusOK,
		usage,
	)
}

func (h *Handler) Create(
	c *gin.Context,
) {
	organizationID, ok :=
		getOrganizationID(c)

	if !ok {
		return
	}

	var request CreateUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "dados inválidos",
			},
		)
		return
	}

	user, err := h.service.Create(
		c.Request.Context(),
		organizationID,
		request,
	)

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(
		http.StatusCreated,
		user,
	)
}

func (h *Handler) UpdateStatus(
	c *gin.Context,
) {
	organizationID, ok :=
		getOrganizationID(c)

	if !ok {
		return
	}

	userID := c.Param("id")

	if userID == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "id do usuário inválido",
			},
		)
		return
	}

	var request UpdateUserStatusRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "dados inválidos",
			},
		)
		return
	}

	user, err := h.service.UpdateStatus(
		c.Request.Context(),
		userID,
		organizationID,
		request.Status,
	)

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(
		http.StatusOK,
		user,
	)
}

func getOrganizationID(
	c *gin.Context,
) (string, bool) {
	value, exists := c.Get(
		shared.ContextOrganizationIDKey,
	)

	if !exists {
		c.JSON(
			http.StatusForbidden,
			gin.H{
				"error": "usuário não possui organização",
			},
		)

		return "", false
	}

	organizationID, ok := value.(string)

	if !ok || organizationID == "" {
		c.JSON(
			http.StatusForbidden,
			gin.H{
				"error": "organização inválida",
			},
		)

		return "", false
	}

	return organizationID, true
}

func handleError(
	c *gin.Context,
	err error,
) {
	switch {
	case errors.Is(
		err,
		ErrLicenseLimitReached,
	):
		c.JSON(
			http.StatusConflict,
			gin.H{
				"error": err.Error(),
				"code":  "LICENSE_LIMIT_REACHED",
			},
		)

	case errors.Is(
		err,
		ErrEmailAlreadyExists,
	):
		c.JSON(
			http.StatusConflict,
			gin.H{
				"error": err.Error(),
				"code":  "EMAIL_ALREADY_EXISTS",
			},
		)

	case errors.Is(
		err,
		ErrUserNotFound,
	):
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": err.Error(),
			},
		)

	case errors.Is(
		err,
		ErrActiveContractNotFound,
	):
		c.JSON(
			http.StatusForbidden,
			gin.H{
				"error": err.Error(),
				"code":  "ACTIVE_CONTRACT_NOT_FOUND",
			},
		)

	case errors.Is(
		err,
		ErrInvalidRole,
	),
		errors.Is(
			err,
			ErrInvalidStatus,
		),
		errors.Is(
			err,
			ErrInvalidPassword,
		),
		errors.Is(
			err,
			ErrInvalidName,
		),
		errors.Is(
			err,
			ErrInvalidEmail,
		):
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

	default:
		log.Printf(
			"[users] erro interno: %v",
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
