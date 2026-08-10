package organizations

import "time"

type Organization struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Document         *string   `json:"document,omitempty"`
	OrganizationType string    `json:"organization_type"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
