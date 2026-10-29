package occurrences

import "time"

const (
	StatusPending      = "pending"
	StatusAcknowledged = "acknowledged"
	StatusAddressed    = "addressed"
	StatusArchived     = "archived"
)

type UpdateOccurrenceStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type OccurrenceStatusHistory struct {
	ID              string    `json:"id"`
	OccurrenceID    string    `json:"occurrence_id"`
	ChangedByUserID *string   `json:"changed_by_user_id,omitempty"`
	PreviousStatus  string    `json:"previous_status"`
	NewStatus       string    `json:"new_status"`
	CreatedAt       time.Time `json:"created_at"`
}

func isValidOccurrenceStatus(status string) bool {
	switch status {
	case StatusPending,
		StatusAcknowledged,
		StatusAddressed,
		StatusArchived:
		return true
	default:
		return false
	}
}
