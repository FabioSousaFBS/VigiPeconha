package healthfacilities

import "time"

type Facility struct {
	ID                 string     `json:"id"`
	Name               string     `json:"name"`
	CNES               *string    `json:"cnes,omitempty"`
	UnitType           string     `json:"unit_type"`
	Latitude           float64    `json:"latitude"`
	Longitude          float64    `json:"longitude"`
	Address            *string    `json:"address,omitempty"`
	City               string     `json:"city"`
	State              string     `json:"state"`
	Phone              *string    `json:"phone,omitempty"`
	Is24Hours          bool       `json:"is_24_hours"`
	TreatsEnvenomation *bool      `json:"treats_envenomation"`
	HasAntivenom       *bool      `json:"has_antivenom"`
	LastVerifiedAt     *time.Time `json:"last_verified_at,omitempty"`
	IsActive           bool       `json:"is_active"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DistanceMeters     *float64   `json:"distance_meters,omitempty"`
}

type SaveRequest struct {
	Name               string     `json:"name" binding:"required"`
	CNES               *string    `json:"cnes"`
	UnitType           string     `json:"unit_type" binding:"required"`
	Latitude           *float64   `json:"latitude" binding:"required"`
	Longitude          *float64   `json:"longitude" binding:"required"`
	Address            *string    `json:"address"`
	City               string     `json:"city" binding:"required"`
	State              string     `json:"state" binding:"required"`
	Phone              *string    `json:"phone"`
	Is24Hours          bool       `json:"is_24_hours"`
	TreatsEnvenomation *bool      `json:"treats_envenomation"`
	HasAntivenom       *bool      `json:"has_antivenom"`
	LastVerifiedAt     *time.Time `json:"last_verified_at"`
	IsActive           *bool      `json:"is_active"`
}

type ListFilter struct {
	State        string
	City         string
	Latitude     *float64
	Longitude    *float64
	RadiusMeters float64
	Limit        int
}
