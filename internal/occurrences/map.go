package occurrences

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MapPoint deliberately excludes reporter details and full descriptions.
type MapPoint struct {
	ID             string    `json:"id"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	AnimalType     string    `json:"animal_type"`
	OccurrenceType string    `json:"occurrence_type"`
	Status         string    `json:"status"`
	OccurredAt     time.Time `json:"occurred_at"`
}

type MapBounds struct{ MinLat, MaxLat, MinLon, MaxLon float64 }

func parseMapBounds(c *gin.Context) (MapBounds, error) {
	keys := []string{"min_lat", "max_lat", "min_lon", "max_lon"}
	vals := make([]float64, 4)
	for i, key := range keys {
		raw := c.Query(key)
		if raw == "" {
			return MapBounds{}, errors.New("limites do mapa são obrigatórios")
		}
		n, err := strconv.ParseFloat(raw, 64)
		if err != nil || math.IsNaN(n) || math.IsInf(n, 0) {
			return MapBounds{}, fmt.Errorf("%s inválido", key)
		}
		vals[i] = n
	}
	b := MapBounds{vals[0], vals[1], vals[2], vals[3]}
	if b.MinLat < -90 || b.MaxLat > 90 || b.MinLon < -180 || b.MaxLon > 180 || b.MinLat >= b.MaxLat || b.MinLon >= b.MaxLon {
		return b, errors.New("limites do mapa inválidos; longitude deve crescer da esquerda para a direita")
	}
	return b, nil
}

// RegisterMapRoutes registers public, read-only, privacy-safe map endpoints.
func RegisterMapRoutes(router *gin.Engine, db *pgxpool.Pool) {
	router.GET("/public/map/occurrences", func(c *gin.Context) {
		b, err := parseMapBounds(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		limit := 500
		if v := c.Query("limit"); v != "" {
			n, e := strconv.Atoi(v)
			if e != nil || n < 1 || n > 1000 {
				c.JSON(400, gin.H{"error": "limit deve estar entre 1 e 1000"})
				return
			}
			limit = n
		}
		status := strings.ToLower(strings.TrimSpace(c.DefaultQuery("status", "active")))
		switch status {
		case "active", "all", "pending", "acknowledged", "addressed", "archived":
		default:
			c.JSON(400, gin.H{"error": "status inválido"})
			return
		}
		animal := strings.TrimSpace(c.Query("animal_type"))
		var from, to *time.Time
		for _, item := range []struct {
			key    string
			target **time.Time
		}{{"from", &from}, {"to", &to}} {
			if v := c.Query(item.key); v != "" {
				t, e := time.Parse(time.RFC3339, v)
				if e != nil {
					c.JSON(400, gin.H{"error": item.key + " deve estar em RFC3339"})
					return
				}
				*item.target = &t
			}
		}
		if from != nil && to != nil && from.After(*to) {
			c.JSON(400, gin.H{"error": "from deve ser anterior a to"})
			return
		}
		items, err := queryOccurrenceMap(c.Request.Context(), db, b, status, animal, from, to, limit)
		if err != nil {
			c.JSON(500, gin.H{"error": "erro ao consultar ocorrências do mapa"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items, "count": len(items), "limit": limit})
	})
}

func queryOccurrenceMap(ctx context.Context, db *pgxpool.Pool, b MapBounds, status, animal string, from, to *time.Time, limit int) ([]MapPoint, error) {
	const query = `SELECT id,ST_Y(location::geometry),ST_X(location::geometry),animal_type,occurrence_type,status,occurred_at
 FROM occurrences WHERE location::geometry && ST_MakeEnvelope($1,$2,$3,$4,4326)
 AND ($5='all' OR ($5='active' AND status<>'archived') OR status=$5)
 AND ($6='' OR LOWER(animal_type)=LOWER($6))
 AND ($7::timestamptz IS NULL OR occurred_at >= $7)
 AND ($8::timestamptz IS NULL OR occurred_at <= $8)
 ORDER BY occurred_at DESC,id LIMIT $9`
	rows, err := db.Query(ctx, query, b.MinLon, b.MinLat, b.MaxLon, b.MaxLat, status, animal, from, to, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]MapPoint, 0)
	for rows.Next() {
		var p MapPoint
		if err := rows.Scan(&p.ID, &p.Latitude, &p.Longitude, &p.AnimalType, &p.OccurrenceType, &p.Status, &p.OccurredAt); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, rows.Err()
}
