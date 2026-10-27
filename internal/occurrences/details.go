package occurrences

type OccurrenceDetailsResponse struct {
	Occurrence
	Photos []OccurrencePhoto `json:"photos"`
}
