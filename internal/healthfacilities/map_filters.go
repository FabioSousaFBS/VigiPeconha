package healthfacilities

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MapFacility is intentionally a compact public projection.
type MapFacility struct {
	ID                 string  `json:"id"`
	Name               string  `json:"name"`
	UnitType           string  `json:"unit_type"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	City               string  `json:"city"`
	State              string  `json:"state"`
	Is24Hours          bool    `json:"is_24_hours"`
	TreatsEnvenomation *bool   `json:"treats_envenomation"`
	HasAntivenom       *bool   `json:"has_antivenom"`
}

// RegisterMapRoutes provides a public bbox endpoint using the same filters as List.
func RegisterMapRoutes(router *gin.Engine, db *pgxpool.Pool) {
	router.GET("/public/map/health-facilities", func(c *gin.Context) {
		keys := []string{"min_lat", "max_lat", "min_lon", "max_lon"}
		v := make([]float64, 4)
		for i, k := range keys {
			raw := c.Query(k)
			if raw == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": k + " obrigatório"})
				return
			}
			n, e := strconv.ParseFloat(raw, 64)
			if e != nil || math.IsNaN(n) || math.IsInf(n, 0) {
				c.JSON(http.StatusBadRequest, gin.H{"error": k + " inválido"})
				return
			}
			v[i] = n
		}
		if v[0] < -90 || v[1] > 90 || v[2] < -180 || v[3] > 180 || v[0] >= v[1] || v[2] >= v[3] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bounding box inválida"})
			return
		}
		limit := 500
		if raw := c.Query("limit"); raw != "" {
			n, e := strconv.Atoi(raw)
			if e != nil || n < 1 || n > 1000 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "limit inválido"})
				return
			}
			limit = n
		}
		f, err := parseFacilityFilters(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		items, err := queryFacilityMap(c.Request.Context(), db, v, f, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao consultar unidades do mapa"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items, "count": len(items), "limit": limit})
	})
}

func queryFacilityMap(ctx context.Context, db *pgxpool.Pool, b []float64, f ListFilter, limit int) ([]MapFacility, error) {
	const q = `SELECT id,name,unit_type,ST_Y(location::geometry),ST_X(location::geometry),city,state,is_24_hours,treats_envenomation,has_antivenom
 FROM health_facilities WHERE is_active=TRUE AND location::geometry && ST_MakeEnvelope($1,$2,$3,$4,4326)
 AND ($5='' OR state=$5) AND ($6='' OR LOWER(city)=LOWER($6))
 AND ($7='' OR unit_type=$7) AND ($8::boolean IS NULL OR treats_envenomation=$8)
 AND ($9::boolean IS NULL OR has_antivenom=$9) ORDER BY name,id LIMIT $10`
	rows, err := db.Query(ctx, q, b[2], b[0], b[3], b[1], strings.ToUpper(strings.TrimSpace(f.State)), strings.TrimSpace(f.City), f.UnitType, f.TreatsEnvenomation, f.HasAntivenom, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]MapFacility, 0)
	for rows.Next() {
		var x MapFacility
		if err := rows.Scan(&x.ID, &x.Name, &x.UnitType, &x.Latitude, &x.Longitude, &x.City, &x.State, &x.Is24Hours, &x.TreatsEnvenomation, &x.HasAntivenom); err != nil {
			return nil, err
		}
		items = append(items, x)
	}
	return items, rows.Err()
}
