package occurrences

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/FabioSousaFBS/vigipeconha-api/internal/storage"
	"github.com/google/uuid"
)

type Service interface {
	CreatePublic(
		ctx context.Context,
		request CreatePublicOccurrenceRequest,
	) (*Occurrence, error)

	UploadPhoto(
		ctx context.Context,
		occurrenceID string,
		size int64,
		content io.Reader,
	) (*OccurrencePhoto, error)

	List(
		ctx context.Context,
		filter ListOccurrencesFilter,
	) (*ListOccurrencesResponse, error)
}

type OccurrenceService struct {
	repository Repository
	storage    storage.Storage
}

func NewService(
	repository Repository,
	fileStorage storage.Storage,
) Service {
	return &OccurrenceService{
		repository: repository,
		storage:    fileStorage,
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

func (s *OccurrenceService) UploadPhoto(
	ctx context.Context,
	occurrenceID string,
	size int64,
	content io.Reader,
) (*OccurrencePhoto, error) {
	const maxFileSize int64 = 5 * 1024 * 1024

	if size <= 0 {
		return nil, ErrEmptyPhoto
	}

	if size > maxFileSize {
		return nil, ErrPhotoTooLarge
	}

	exists, err := s.repository.Exists(
		ctx,
		occurrenceID,
	)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrOccurrenceNotFound
	}

	data, err := io.ReadAll(
		io.LimitReader(content, maxFileSize+1),
	)
	if err != nil {
		return nil, ErrInvalidPhoto
	}

	if int64(len(data)) > maxFileSize {
		return nil, ErrPhotoTooLarge
	}

	if len(data) == 0 {
		return nil, ErrEmptyPhoto
	}

	contentType := http.DetectContentType(data)

	if !isAllowedImageType(contentType) {
		return nil, ErrInvalidPhotoType
	}

	extension := extensionForContentType(contentType)

	key := fmt.Sprintf(
		"occurrences/%s/%s%s",
		occurrenceID,
		uuid.NewString(),
		extension,
	)

	err = s.storage.Upload(
		ctx,
		key,
		bytes.NewReader(data),
		contentType,
	)
	if err != nil {
		return nil, err
	}

	photo := OccurrencePhoto{
		OccurrenceID: occurrenceID,
		StorageKey:   key,
		PhotoURL:     nil,
		ContentType:  contentType,
		FileSize:     int64(len(data)),
	}

	created, err := s.repository.CreatePhoto(
		ctx,
		photo,
	)
	if err != nil {
		_ = s.storage.Delete(ctx, key)

		return nil, err
	}

	return created, nil
}

func isAllowedImageType(
	contentType string,
) bool {
	switch contentType {
	case "image/jpeg",
		"image/png",
		"image/webp":
		return true

	default:
		return false
	}
}

func extensionForContentType(
	contentType string,
) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"

	case "image/png":
		return ".png"

	case "image/webp":
		return ".webp"

	default:
		return ""
	}
}

func (s *OccurrenceService) List(
	ctx context.Context,
	filter ListOccurrencesFilter,
) (*ListOccurrencesResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}

	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	filter.Status = normalizeOptionalString(
		filter.Status,
	)

	filter.OccurrenceType = normalizeOptionalString(
		filter.OccurrenceType,
	)

	filter.AnimalType = normalizeOptionalString(
		filter.AnimalType,
	)

	filter.State = normalizeOptionalString(
		filter.State,
	)

	filter.City = normalizeOptionalString(
		filter.City,
	)

	if filter.Status != nil {
		normalized :=
			strings.ToLower(*filter.Status)

		filter.Status = &normalized
	}

	if filter.OccurrenceType != nil {
		normalized :=
			strings.ToLower(*filter.OccurrenceType)

		if !isValidOccurrenceType(normalized) {
			return nil, ErrInvalidOccurrenceType
		}

		filter.OccurrenceType = &normalized
	}

	if filter.State != nil {
		normalized :=
			strings.ToUpper(*filter.State)

		if len(normalized) != 2 {
			return nil, ErrInvalidState
		}

		filter.State = &normalized
	}

	items, total, err := s.repository.List(
		ctx,
		filter,
	)
	if err != nil {
		return nil, err
	}

	totalPages := 0

	if total > 0 {
		totalPages =
			int(
				(total + int64(filter.PageSize) - 1) /
					int64(filter.PageSize),
			)
	}

	return &ListOccurrencesResponse{
		Items:      items,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}
