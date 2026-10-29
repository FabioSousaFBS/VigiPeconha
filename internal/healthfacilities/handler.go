package healthfacilities

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct{ service Service }

func NewHandler(s Service) *Handler { return &Handler{service: s} }
func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidInput), errors.Is(err, ErrInvalidID):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, ErrDuplicateCNES):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		log.Printf("[healthfacilities] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro interno do servidor"})
	}
}
func (h *Handler) Create(c *gin.Context) {
	var r SaveRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(400, gin.H{"error": "dados inválidos"})
		return
	}
	f, err := h.service.Create(c.Request.Context(), r)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, f)
}
func (h *Handler) Update(c *gin.Context) {
	var r SaveRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(400, gin.H{"error": "dados inválidos"})
		return
	}
	f, err := h.service.Update(c.Request.Context(), c.Param("id"), r)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(200, f)
}
func (h *Handler) GetByID(c *gin.Context) {
	f, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(200, f)
}
func (h *Handler) List(c *gin.Context) {
	f := ListFilter{State: c.Query("state"), City: c.Query("city"), Limit: 50}
	if v := c.Query("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			c.JSON(400, gin.H{"error": "limit inválido"})
			return
		}
		f.Limit = n
	}
	if c.FullPath() == "/public/health-facilities/nearby" {
		lat, err1 := strconv.ParseFloat(c.Query("latitude"), 64)
		lon, err2 := strconv.ParseFloat(c.Query("longitude"), 64)
		if err1 != nil || err2 != nil {
			c.JSON(400, gin.H{"error": "latitude e longitude obrigatórias"})
			return
		}
		radius := 10000.0
		if v := c.Query("radius_meters"); v != "" {
			var err error
			radius, err = strconv.ParseFloat(v, 64)
			if err != nil {
				c.JSON(400, gin.H{"error": "radius_meters inválido"})
				return
			}
		}
		f.Latitude = &lat
		f.Longitude = &lon
		f.RadiusMeters = radius
	}
	items, err := h.service.List(c.Request.Context(), f)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(200, gin.H{"items": items, "total": len(items)})
}
