package occurrences

import (
	"context"
	"strings"
)

type Service interface {
	CreatePublic(
		ctx context.Context,
		request CreatePublicOccurrenceRequest,
	) (*Occurrence, error)
}

type OccurrenceService struct {
	repository Repository
}

func NewService(
	repository Repository,
) Service {
	return &OccurrenceService{
		repository: repository,
	}
}

func (s *OccurrenceService) CreatePublic(
	ctx context.Context,
	request CreatePublicOccurrenceRequest,
) (*Occurrence, error) {
	occurrenceType := strings.ToLower(
		strings.TrimSpace(
			request.OccurrenceType,
		),
	)

	if !isValidOccurrenceType(
		occurrenceType,
	) {
		return nil, ErrInvalidOccurrenceType
	}

	animalType := strings.TrimSpace(
		request.AnimalType,
	)

	if animalType == "" {
		return nil, ErrAnimalTypeRequired
	}

	if request.Latitude == nil ||
		request.Longitude == nil {
		return nil, ErrCoordinatesRequired
	}

	if *request.Latitude < -90 ||
		*request.Latitude > 90 {
		return nil, ErrInvalidLatitude
	}

	if *request.Longitude < -180 ||
		*request.Longitude > 180 {
		return nil, ErrInvalidLongitude
	}

	if request.OccurredAt.IsZero() {
		return nil, ErrOccurredAtRequired
	}

	state := normalizeOptionalString(
		request.State,
	)

	if state != nil {
		normalizedState :=
			strings.ToUpper(*state)

		if len(normalizedState) != 2 {
			return nil, ErrInvalidState
		}

		state = &normalizedState
	}

	params := CreateOccurrenceParams{
		OrganizationID:  nil,
		CreatedByUserID: nil,

		Source: "mobile_public",

		OccurrenceType: occurrenceType,

		AnimalType: animalType,

		Species: normalizeOptionalString(
			request.Species,
		),

		Description: normalizeOptionalString(
			request.Description,
		),

		ReporterName: normalizeOptionalString(
			request.ReporterName,
		),

		ReporterPhone: normalizeOptionalString(
			request.ReporterPhone,
		),

		Latitude: *request.Latitude,

		Longitude: *request.Longitude,

		Address: normalizeOptionalString(
			request.Address,
		),

		Neighborhood: normalizeOptionalString(
			request.Neighborhood,
		),

		City: normalizeOptionalString(
			request.City,
		),

		State: state,

		Status: "pending",

		OccurredAt: request.OccurredAt,
	}

	return s.repository.Create(
		ctx,
		params,
	)
}

func isValidOccurrenceType(
	value string,
) bool {
	switch value {
	case "sighting",
		"accident",
		"capture":
		return true

	default:
		return false
	}
}

func normalizeOptionalString(
	value *string,
) *string {
	if value == nil {
		return nil
	}

	normalized :=
		strings.TrimSpace(*value)

	if normalized == "" {
		return nil
	}

	return &normalized
}
