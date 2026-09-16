package occurrences

import "time"

type OccurrencePhoto struct {
	ID           string    `json:"id"`
	OccurrenceID string    `json:"occurrence_id"`
	StorageKey   string    `json:"storage_key"`
	PhotoURL     *string   `json:"photo_url,omitempty"`
	ContentType  string    `json:"content_type"`
	FileSize     int64     `json:"file_size"`
	CreatedAt    time.Time `json:"created_at"`
}
