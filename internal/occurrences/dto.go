package occurrences

import "time"

type CreatePublicOccurrenceRequest struct {
	OccurrenceType string  `json:"occurrence_type"`
	AnimalType     string  `json:"animal_type"`
	Species        *string `json:"species"`
	Description    *string `json:"description"`

	ReporterName  *string `json:"reporter_name"`
	ReporterPhone *string `json:"reporter_phone"`

	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`

	Address      *string `json:"address"`
	Neighborhood *string `json:"neighborhood"`
	City         *string `json:"city"`
	State        *string `json:"state"`

	OccurredAt time.Time `json:"occurred_at"`
}

type CreateOccurrenceParams struct {
	OrganizationID  *string
	CreatedByUserID *string

	Source         string
	OccurrenceType string
	AnimalType     string
	Species        *string
	Description    *string

	ReporterName  *string
	ReporterPhone *string

	Latitude  float64
	Longitude float64

	Address      *string
	Neighborhood *string
	City         *string
	State        *string

	Status     string
	OccurredAt time.Time
}

type ListOccurrencesFilter struct {
	Page           int
	PageSize       int
	Status         *string
	OccurrenceType *string
	AnimalType     *string
	State          *string
	City           *string
}

type ListOccurrencesResponse struct {
	Items      []Occurrence `json:"items"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	Total      int64        `json:"total"`
	TotalPages int          `json:"total_pages"`
}