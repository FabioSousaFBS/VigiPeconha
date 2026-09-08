package occurrences

import "time"

type Occurrence struct {
	ID              string  `json:"id"`
	OrganizationID  *string `json:"organization_id,omitempty"`
	CreatedByUserID *string `json:"created_by_user_id,omitempty"`

	Source         string  `json:"source"`
	OccurrenceType string  `json:"occurrence_type"`
	AnimalType     string  `json:"animal_type"`
	Species        *string `json:"species,omitempty"`
	Description    *string `json:"description,omitempty"`

	ReporterName  *string `json:"reporter_name,omitempty"`
	ReporterPhone *string `json:"reporter_phone,omitempty"`

	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	Address      *string `json:"address,omitempty"`
	Neighborhood *string `json:"neighborhood,omitempty"`
	City         *string `json:"city,omitempty"`
	State        *string `json:"state,omitempty"`

	Status string `json:"status"`

	OccurredAt time.Time `json:"occurred_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
