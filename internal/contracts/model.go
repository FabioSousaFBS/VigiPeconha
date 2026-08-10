package contracts

import "time"

type Contract struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	UserLimit      int       `json:"user_limit"`
	Status         string    `json:"status"`
	PaymentStatus  string    `json:"payment_status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
