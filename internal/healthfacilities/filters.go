package healthfacilities

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func parseFacilityFilters(c *gin.Context) (ListFilter, error) {
	f := ListFilter{
		City:     strings.TrimSpace(c.Query("city")),
		State:    strings.ToUpper(strings.TrimSpace(c.Query("state"))),
		UnitType: strings.ToLower(strings.TrimSpace(c.Query("unit_type"))),
	}

	if raw := c.Query("treats_envenomation"); raw != "" {
		value, err := parseBoolFilter(raw)
		if err != nil {
			return f, fmt.Errorf("treats_envenomation inválido: use true ou false")
		}
		f.TreatsEnvenomation = &value
	}

	if raw := c.Query("has_antivenom"); raw != "" {
		value, err := parseBoolFilter(raw)
		if err != nil {
			return f, fmt.Errorf("has_antivenom inválido: use true ou false")
		}
		f.HasAntivenom = &value
	}

	return f, nil
}

func parseBoolFilter(raw string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("valor booleano inválido")
	}
}
