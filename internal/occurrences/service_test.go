package occurrences

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"
)

type mockRepository struct {
	exists    bool
	existsErr error

	createCalled bool
	createParams CreateOccurrenceParams
	createResult *Occurrence
	createErr    error

	createPhotoCalled bool
	createPhotoErr    error

	listCalled bool
	listFilter ListOccurrencesFilter
	listResult []Occurrence
	listTotal  int64
	listErr    error
}

func (m *mockRepository) Create(
	ctx context.Context,
	params CreateOccurrenceParams,
) (*Occurrence, error) {
	m.createCalled = true
	m.createParams = params

	if m.createErr != nil {
		return nil, m.createErr
	}

	if m.createResult != nil {
		return m.createResult, nil
	}

	return &Occurrence{}, nil
}

func (m *mockRepository) CreatePhoto(
	ctx context.Context,
	photo OccurrencePhoto,
) (*OccurrencePhoto, error) {
	m.createPhotoCalled = true

	if m.createPhotoErr != nil {
		return nil, m.createPhotoErr
	}

	return &photo, nil
}

func (m *mockRepository) Exists(
	ctx context.Context,
	occurrenceID string,
) (bool, error) {
	if m.existsErr != nil {
		return false, m.existsErr
	}

	return m.exists, nil
}

type mockStorage struct {
	uploadCalled bool
	deleteCalled bool

	uploadedKey string
	deletedKey  string

	uploadErr error
}

func (m *mockStorage) Upload(
	ctx context.Context,
	key string,
	body io.Reader,
	contentType string,
) error {
	m.uploadCalled = true
	m.uploadedKey = key

	return m.uploadErr
}

func (m *mockStorage) Delete(
	ctx context.Context,
	key string,
) error {
	m.deleteCalled = true
	m.deletedKey = key

	return nil
}

func (m *mockRepository) List(
	ctx context.Context,
	filter ListOccurrencesFilter,
) ([]Occurrence, int64, error) {
	m.listCalled = true
	m.listFilter = filter

	if m.listErr != nil {
		return nil, 0, m.listErr
	}

	return m.listResult, m.listTotal, nil
}

func TestCreatePublicSuccess(t *testing.T) {
	latitude := -22.2100
	longitude := -49.6500
	state := "sp"

	expectedOccurrence := &Occurrence{
		ID:             "occurrence-id",
		Source:         "mobile_public",
		OccurrenceType: "sighting",
		AnimalType:     "snake",
		Latitude:       latitude,
		Longitude:      longitude,
		Status:         "pending",
	}

	repository := &mockRepository{
		createResult: expectedOccurrence,
	}

	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	request := CreatePublicOccurrenceRequest{
		OccurrenceType: "sighting",
		AnimalType:     "snake",
		Latitude:       &latitude,
		Longitude:      &longitude,
		State:          &state,
		OccurredAt:     time.Now(),
	}

	result, err := service.CreatePublic(
		context.Background(),
		request,
	)

	if err != nil {
		t.Fatalf(
			"não esperava erro, recebeu %v",
			err,
		)
	}

	if result == nil {
		t.Fatal("esperava uma ocorrência, recebeu nil")
	}

	if !repository.createCalled {
		t.Error("Repository.Create deveria ser chamado")
	}

	if result.ID != expectedOccurrence.ID {
		t.Errorf(
			"esperava ID %s, recebeu %s",
			expectedOccurrence.ID,
			result.ID,
		)
	}

	if repository.createParams.Source != "mobile_public" {
		t.Errorf(
			"esperava source mobile_public, recebeu %s",
			repository.createParams.Source,
		)
	}

	if repository.createParams.Status != "pending" {
		t.Errorf(
			"esperava status pending, recebeu %s",
			repository.createParams.Status,
		)
	}

	if repository.createParams.OrganizationID != nil {
		t.Error("OrganizationID deveria ser nil")
	}

	if repository.createParams.CreatedByUserID != nil {
		t.Error("CreatedByUserID deveria ser nil")
	}

	if repository.createParams.State == nil {
		t.Fatal("State não deveria ser nil")
	}

	if *repository.createParams.State != "SP" {
		t.Errorf(
			"esperava estado SP, recebeu %s",
			*repository.createParams.State,
		)
	}
}

func TestCreatePublicNormalizesFields(t *testing.T) {
	latitude := -22.2100
	longitude := -49.6500

	species := "  Bothrops jararaca  "
	description := "  encontrada no quintal  "
	reporterName := "  Fabio  "
	reporterPhone := "  14999999999  "
	address := "  Rua Teste  "
	neighborhood := "  Centro  "
	city := "  Garça  "
	state := "  sp  "

	repository := &mockRepository{}
	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	request := CreatePublicOccurrenceRequest{
		OccurrenceType: "  SIGHTING  ",
		AnimalType:     "  snake  ",
		Species:        &species,
		Description:    &description,
		ReporterName:   &reporterName,
		ReporterPhone:  &reporterPhone,
		Latitude:       &latitude,
		Longitude:      &longitude,
		Address:        &address,
		Neighborhood:   &neighborhood,
		City:           &city,
		State:          &state,
		OccurredAt:     time.Now(),
	}

	_, err := service.CreatePublic(
		context.Background(),
		request,
	)

	if err != nil {
		t.Fatalf(
			"não esperava erro, recebeu %v",
			err,
		)
	}

	params := repository.createParams

	if params.OccurrenceType != "sighting" {
		t.Errorf(
			"esperava sighting, recebeu %s",
			params.OccurrenceType,
		)
	}

	if params.AnimalType != "snake" {
		t.Errorf(
			"esperava snake, recebeu %s",
			params.AnimalType,
		)
	}

	assertStringPointerEquals(
		t,
		params.Species,
		"Bothrops jararaca",
		"Species",
	)

	assertStringPointerEquals(
		t,
		params.Description,
		"encontrada no quintal",
		"Description",
	)

	assertStringPointerEquals(
		t,
		params.ReporterName,
		"Fabio",
		"ReporterName",
	)

	assertStringPointerEquals(
		t,
		params.ReporterPhone,
		"14999999999",
		"ReporterPhone",
	)

	assertStringPointerEquals(
		t,
		params.Address,
		"Rua Teste",
		"Address",
	)

	assertStringPointerEquals(
		t,
		params.Neighborhood,
		"Centro",
		"Neighborhood",
	)

	assertStringPointerEquals(
		t,
		params.City,
		"Garça",
		"City",
	)

	assertStringPointerEquals(
		t,
		params.State,
		"SP",
		"State",
	)
}

func TestCreatePublicEmptyOptionalFieldsBecomeNil(
	t *testing.T,
) {
	latitude := -22.2100
	longitude := -49.6500

	empty := "   "

	repository := &mockRepository{}
	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	request := CreatePublicOccurrenceRequest{
		OccurrenceType: "sighting",
		AnimalType:     "snake",
		Species:        &empty,
		Description:    &empty,
		ReporterName:   &empty,
		ReporterPhone:  &empty,
		Latitude:       &latitude,
		Longitude:      &longitude,
		Address:        &empty,
		Neighborhood:   &empty,
		City:           &empty,
		State:          nil,
		OccurredAt:     time.Now(),
	}

	_, err := service.CreatePublic(
		context.Background(),
		request,
	)

	if err != nil {
		t.Fatalf(
			"não esperava erro, recebeu %v",
			err,
		)
	}

	params := repository.createParams

	if params.Species != nil {
		t.Error("Species deveria ser nil")
	}

	if params.Description != nil {
		t.Error("Description deveria ser nil")
	}

	if params.ReporterName != nil {
		t.Error("ReporterName deveria ser nil")
	}

	if params.ReporterPhone != nil {
		t.Error("ReporterPhone deveria ser nil")
	}

	if params.Address != nil {
		t.Error("Address deveria ser nil")
	}

	if params.Neighborhood != nil {
		t.Error("Neighborhood deveria ser nil")
	}

	if params.City != nil {
		t.Error("City deveria ser nil")
	}

	if params.State != nil {
		t.Error("State deveria ser nil")
	}
}

func TestCreatePublicInvalidOccurrenceType(t *testing.T) {
	latitude := -22.2100
	longitude := -49.6500

	repository := &mockRepository{}
	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	request := CreatePublicOccurrenceRequest{
		OccurrenceType: "invalid",
		AnimalType:     "snake",
		Latitude:       &latitude,
		Longitude:      &longitude,
		OccurredAt:     time.Now(),
	}

	_, err := service.CreatePublic(
		context.Background(),
		request,
	)

	if !errors.Is(err, ErrInvalidOccurrenceType) {
		t.Errorf(
			"esperava ErrInvalidOccurrenceType, recebeu %v",
			err,
		)
	}

	if repository.createCalled {
		t.Error(
			"Repository.Create não deveria ser chamado",
		)
	}
}

func TestCreatePublicAnimalTypeRequired(t *testing.T) {
	latitude := -22.2100
	longitude := -49.6500

	repository := &mockRepository{}
	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	request := CreatePublicOccurrenceRequest{
		OccurrenceType: "sighting",
		AnimalType:     "   ",
		Latitude:       &latitude,
		Longitude:      &longitude,
		OccurredAt:     time.Now(),
	}

	_, err := service.CreatePublic(
		context.Background(),
		request,
	)

	if !errors.Is(err, ErrAnimalTypeRequired) {
		t.Errorf(
			"esperava ErrAnimalTypeRequired, recebeu %v",
			err,
		)
	}

	if repository.createCalled {
		t.Error(
			"Repository.Create não deveria ser chamado",
		)
	}
}

func TestCreatePublicCoordinatesRequired(t *testing.T) {
	repository := &mockRepository{}
	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	request := CreatePublicOccurrenceRequest{
		OccurrenceType: "sighting",
		AnimalType:     "snake",
		Latitude:       nil,
		Longitude:      nil,
		OccurredAt:     time.Now(),
	}

	_, err := service.CreatePublic(
		context.Background(),
		request,
	)

	if !errors.Is(err, ErrCoordinatesRequired) {
		t.Errorf(
			"esperava ErrCoordinatesRequired, recebeu %v",
			err,
		)
	}

	if repository.createCalled {
		t.Error(
			"Repository.Create não deveria ser chamado",
		)
	}
}

func TestCreatePublicInvalidLatitude(t *testing.T) {
	latitude := 91.0
	longitude := -49.6500

	repository := &mockRepository{}
	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	request := CreatePublicOccurrenceRequest{
		OccurrenceType: "sighting",
		AnimalType:     "snake",
		Latitude:       &latitude,
		Longitude:      &longitude,
		OccurredAt:     time.Now(),
	}

	_, err := service.CreatePublic(
		context.Background(),
		request,
	)

	if !errors.Is(err, ErrInvalidLatitude) {
		t.Errorf(
			"esperava ErrInvalidLatitude, recebeu %v",
			err,
		)
	}

	if repository.createCalled {
		t.Error(
			"Repository.Create não deveria ser chamado",
		)
	}
}

func TestCreatePublicInvalidLongitude(t *testing.T) {
	latitude := -22.2100
	longitude := 181.0

	repository := &mockRepository{}
	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	request := CreatePublicOccurrenceRequest{
		OccurrenceType: "sighting",
		AnimalType:     "snake",
		Latitude:       &latitude,
		Longitude:      &longitude,
		OccurredAt:     time.Now(),
	}

	_, err := service.CreatePublic(
		context.Background(),
		request,
	)

	if !errors.Is(err, ErrInvalidLongitude) {
		t.Errorf(
			"esperava ErrInvalidLongitude, recebeu %v",
			err,
		)
	}

	if repository.createCalled {
		t.Error(
			"Repository.Create não deveria ser chamado",
		)
	}
}

func TestCreatePublicOccurredAtRequired(t *testing.T) {
	latitude := -22.2100
	longitude := -49.6500

	repository := &mockRepository{}
	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	request := CreatePublicOccurrenceRequest{
		OccurrenceType: "sighting",
		AnimalType:     "snake",
		Latitude:       &latitude,
		Longitude:      &longitude,
	}

	_, err := service.CreatePublic(
		context.Background(),
		request,
	)

	if !errors.Is(err, ErrOccurredAtRequired) {
		t.Errorf(
			"esperava ErrOccurredAtRequired, recebeu %v",
			err,
		)
	}

	if repository.createCalled {
		t.Error(
			"Repository.Create não deveria ser chamado",
		)
	}
}

func TestCreatePublicInvalidState(t *testing.T) {
	latitude := -22.2100
	longitude := -49.6500
	state := "SPP"

	repository := &mockRepository{}
	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	request := CreatePublicOccurrenceRequest{
		OccurrenceType: "sighting",
		AnimalType:     "snake",
		Latitude:       &latitude,
		Longitude:      &longitude,
		State:          &state,
		OccurredAt:     time.Now(),
	}

	_, err := service.CreatePublic(
		context.Background(),
		request,
	)

	if !errors.Is(err, ErrInvalidState) {
		t.Errorf(
			"esperava ErrInvalidState, recebeu %v",
			err,
		)
	}

	if repository.createCalled {
		t.Error(
			"Repository.Create não deveria ser chamado",
		)
	}
}

func TestCreatePublicRepositoryError(t *testing.T) {
	latitude := -22.2100
	longitude := -49.6500

	expectedError := errors.New(
		"erro simulado no repository",
	)

	repository := &mockRepository{
		createErr: expectedError,
	}

	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	request := CreatePublicOccurrenceRequest{
		OccurrenceType: "sighting",
		AnimalType:     "snake",
		Latitude:       &latitude,
		Longitude:      &longitude,
		OccurredAt:     time.Now(),
	}

	result, err := service.CreatePublic(
		context.Background(),
		request,
	)

	if !errors.Is(err, expectedError) {
		t.Errorf(
			"esperava erro do repository, recebeu %v",
			err,
		)
	}

	if result != nil {
		t.Error(
			"resultado deveria ser nil quando Repository.Create falha",
		)
	}

	if !repository.createCalled {
		t.Error(
			"Repository.Create deveria ser chamado",
		)
	}
}

func TestUploadPhotoEmptyFile(t *testing.T) {
	repository := &mockRepository{}
	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	_, err := service.UploadPhoto(
		context.Background(),
		"occurrence-id",
		0,
		bytes.NewReader(nil),
	)

	if !errors.Is(err, ErrEmptyPhoto) {
		t.Errorf(
			"esperava ErrEmptyPhoto, recebeu %v",
			err,
		)
	}
}

func TestUploadPhotoTooLarge(t *testing.T) {
	repository := &mockRepository{}
	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	const fileSize int64 = 5*1024*1024 + 1

	_, err := service.UploadPhoto(
		context.Background(),
		"occurrence-id",
		fileSize,
		bytes.NewReader(nil),
	)

	if !errors.Is(err, ErrPhotoTooLarge) {
		t.Errorf(
			"esperava ErrPhotoTooLarge, recebeu %v",
			err,
		)
	}
}

func TestUploadPhotoInvalidContentType(t *testing.T) {
	repository := &mockRepository{
		exists: true,
	}

	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	data := []byte(
		"este arquivo nao e uma imagem",
	)

	_, err := service.UploadPhoto(
		context.Background(),
		"occurrence-id",
		int64(len(data)),
		bytes.NewReader(data),
	)

	if !errors.Is(err, ErrInvalidPhotoType) {
		t.Errorf(
			"esperava ErrInvalidPhotoType, recebeu %v",
			err,
		)
	}
}

func TestUploadPhotoOccurrenceNotFound(t *testing.T) {
	repository := &mockRepository{
		exists: false,
	}

	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	pngData := validPNGData()

	_, err := service.UploadPhoto(
		context.Background(),
		"occurrence-id",
		int64(len(pngData)),
		bytes.NewReader(pngData),
	)

	if !errors.Is(err, ErrOccurrenceNotFound) {
		t.Errorf(
			"esperava ErrOccurrenceNotFound, recebeu %v",
			err,
		)
	}
}

func TestUploadPhotoRepositoryExistsError(t *testing.T) {
	expectedError := errors.New(
		"erro ao verificar ocorrência",
	)

	repository := &mockRepository{
		existsErr: expectedError,
	}

	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	pngData := validPNGData()

	_, err := service.UploadPhoto(
		context.Background(),
		"occurrence-id",
		int64(len(pngData)),
		bytes.NewReader(pngData),
	)

	if !errors.Is(err, expectedError) {
		t.Errorf(
			"esperava erro do repository, recebeu %v",
			err,
		)
	}

	if storage.uploadCalled {
		t.Error(
			"Storage.Upload não deveria ser chamado",
		)
	}
}

func TestUploadPhotoSuccess(t *testing.T) {
	repository := &mockRepository{
		exists: true,
	}

	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	pngData := validPNGData()

	photo, err := service.UploadPhoto(
		context.Background(),
		"occurrence-id",
		int64(len(pngData)),
		bytes.NewReader(pngData),
	)

	if err != nil {
		t.Fatalf(
			"não esperava erro, recebeu %v",
			err,
		)
	}

	if photo == nil {
		t.Fatal(
			"esperava uma foto, recebeu nil",
		)
	}

	if !storage.uploadCalled {
		t.Error(
			"Storage.Upload deveria ser chamado",
		)
	}

	if !repository.createPhotoCalled {
		t.Error(
			"Repository.CreatePhoto deveria ser chamado",
		)
	}

	if storage.deleteCalled {
		t.Error(
			"Storage.Delete não deveria ser chamado",
		)
	}

	if photo.ContentType != "image/png" {
		t.Errorf(
			"esperava image/png, recebeu %s",
			photo.ContentType,
		)
	}

	if photo.FileSize != int64(len(pngData)) {
		t.Errorf(
			"esperava tamanho %d, recebeu %d",
			len(pngData),
			photo.FileSize,
		)
	}

	if photo.OccurrenceID != "occurrence-id" {
		t.Errorf(
			"esperava occurrence-id, recebeu %s",
			photo.OccurrenceID,
		)
	}

	if storage.uploadedKey == "" {
		t.Error(
			"storage key não deveria estar vazia",
		)
	}
}

func TestUploadPhotoStorageError(t *testing.T) {
	repository := &mockRepository{
		exists: true,
	}

	expectedError := errors.New(
		"erro simulado no R2",
	)

	storage := &mockStorage{
		uploadErr: expectedError,
	}

	service := NewService(
		repository,
		storage,
	)

	pngData := validPNGData()

	_, err := service.UploadPhoto(
		context.Background(),
		"occurrence-id",
		int64(len(pngData)),
		bytes.NewReader(pngData),
	)

	if !errors.Is(err, expectedError) {
		t.Errorf(
			"esperava erro do storage, recebeu %v",
			err,
		)
	}

	if !storage.uploadCalled {
		t.Error(
			"Storage.Upload deveria ser chamado",
		)
	}

	if repository.createPhotoCalled {
		t.Error(
			"Repository.CreatePhoto não deveria ser chamado",
		)
	}
}

func TestUploadPhotoDeletesFileWhenDatabaseFails(
	t *testing.T,
) {
	expectedError := errors.New(
		"erro simulado no banco",
	)

	repository := &mockRepository{
		exists:         true,
		createPhotoErr: expectedError,
	}

	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	pngData := validPNGData()

	_, err := service.UploadPhoto(
		context.Background(),
		"occurrence-id",
		int64(len(pngData)),
		bytes.NewReader(pngData),
	)

	if !errors.Is(err, expectedError) {
		t.Errorf(
			"esperava erro do banco, recebeu %v",
			err,
		)
	}

	if !storage.uploadCalled {
		t.Error(
			"Storage.Upload deveria ser chamado",
		)
	}

	if !repository.createPhotoCalled {
		t.Error(
			"Repository.CreatePhoto deveria ser chamado",
		)
	}

	if !storage.deleteCalled {
		t.Error(
			"Storage.Delete deveria ser chamado",
		)
	}

	if storage.uploadedKey != storage.deletedKey {
		t.Errorf(
			"esperava remover %s, removeu %s",
			storage.uploadedKey,
			storage.deletedKey,
		)
	}
}

func validPNGData() []byte {
	return []byte{
		0x89, 0x50, 0x4E, 0x47,
		0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D,
		0x49, 0x48, 0x44, 0x52,
	}
}

func assertStringPointerEquals(
	t *testing.T,
	actual *string,
	expected string,
	field string,
) {
	t.Helper()

	if actual == nil {
		t.Fatalf(
			"%s não deveria ser nil",
			field,
		)
	}

	if *actual != expected {
		t.Errorf(
			"esperava %s=%s, recebeu %s",
			field,
			expected,
			*actual,
		)
	}
}

func TestListOccurrencesSuccess(t *testing.T) {
	repository := &mockRepository{
		listResult: []Occurrence{
			{
				ID:             "occurrence-1",
				Source:         "mobile_public",
				OccurrenceType: "sighting",
				AnimalType:     "snake",
				Status:         "pending",
			},
			{
				ID:             "occurrence-2",
				Source:         "mobile_public",
				OccurrenceType: "accident",
				AnimalType:     "scorpion",
				Status:         "pending",
			},
		},
		listTotal: 2,
	}

	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	result, err := service.List(
		context.Background(),
		ListOccurrencesFilter{
			Page:     1,
			PageSize: 20,
		},
	)

	if err != nil {
		t.Fatalf(
			"não esperava erro, recebeu %v",
			err,
		)
	}

	if result == nil {
		t.Fatal(
			"esperava resultado, recebeu nil",
		)
	}

	if !repository.listCalled {
		t.Error(
			"Repository.List deveria ser chamado",
		)
	}

	if len(result.Items) != 2 {
		t.Errorf(
			"esperava 2 ocorrências, recebeu %d",
			len(result.Items),
		)
	}

	if result.Total != 2 {
		t.Errorf(
			"esperava total 2, recebeu %d",
			result.Total,
		)
	}

	if result.Page != 1 {
		t.Errorf(
			"esperava page 1, recebeu %d",
			result.Page,
		)
	}

	if result.PageSize != 20 {
		t.Errorf(
			"esperava page_size 20, recebeu %d",
			result.PageSize,
		)
	}

	if result.TotalPages != 1 {
		t.Errorf(
			"esperava total_pages 1, recebeu %d",
			result.TotalPages,
		)
	}
}

func TestListOccurrencesDefaultPagination(t *testing.T) {
	repository := &mockRepository{}

	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	result, err := service.List(
		context.Background(),
		ListOccurrencesFilter{},
	)

	if err != nil {
		t.Fatalf(
			"não esperava erro, recebeu %v",
			err,
		)
	}

	if repository.listFilter.Page != 1 {
		t.Errorf(
			"esperava page 1, recebeu %d",
			repository.listFilter.Page,
		)
	}

	if repository.listFilter.PageSize != 20 {
		t.Errorf(
			"esperava page_size 20, recebeu %d",
			repository.listFilter.PageSize,
		)
	}

	if result.Page != 1 {
		t.Errorf(
			"esperava response page 1, recebeu %d",
			result.Page,
		)
	}

	if result.PageSize != 20 {
		t.Errorf(
			"esperava response page_size 20, recebeu %d",
			result.PageSize,
		)
	}
}

func TestListOccurrencesLimitsPageSize(t *testing.T) {
	repository := &mockRepository{}

	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	result, err := service.List(
		context.Background(),
		ListOccurrencesFilter{
			Page:     1,
			PageSize: 1000,
		},
	)

	if err != nil {
		t.Fatalf(
			"não esperava erro, recebeu %v",
			err,
		)
	}

	if repository.listFilter.PageSize != 100 {
		t.Errorf(
			"esperava page_size limitado a 100, recebeu %d",
			repository.listFilter.PageSize,
		)
	}

	if result.PageSize != 100 {
		t.Errorf(
			"esperava response page_size 100, recebeu %d",
			result.PageSize,
		)
	}
}

func TestListOccurrencesCalculatesTotalPages(t *testing.T) {
	repository := &mockRepository{
		listTotal: 45,
	}

	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	result, err := service.List(
		context.Background(),
		ListOccurrencesFilter{
			Page:     1,
			PageSize: 20,
		},
	)

	if err != nil {
		t.Fatalf(
			"não esperava erro, recebeu %v",
			err,
		)
	}

	if result.Total != 45 {
		t.Errorf(
			"esperava total 45, recebeu %d",
			result.Total,
		)
	}

	if result.TotalPages != 3 {
		t.Errorf(
			"esperava total_pages 3, recebeu %d",
			result.TotalPages,
		)
	}
}

func TestListOccurrencesNormalizesFilters(t *testing.T) {
	status := "  PENDING  "
	occurrenceType := "  SIGHTING  "
	animalType := "  snake  "
	state := "  sp  "
	city := "  Garça  "

	repository := &mockRepository{}

	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	_, err := service.List(
		context.Background(),
		ListOccurrencesFilter{
			Page:           1,
			PageSize:       20,
			Status:         &status,
			OccurrenceType: &occurrenceType,
			AnimalType:     &animalType,
			State:          &state,
			City:           &city,
		},
	)

	if err != nil {
		t.Fatalf(
			"não esperava erro, recebeu %v",
			err,
		)
	}

	filter := repository.listFilter

	assertStringPointerEquals(
		t,
		filter.Status,
		"pending",
		"Status",
	)

	assertStringPointerEquals(
		t,
		filter.OccurrenceType,
		"sighting",
		"OccurrenceType",
	)

	assertStringPointerEquals(
		t,
		filter.AnimalType,
		"snake",
		"AnimalType",
	)

	assertStringPointerEquals(
		t,
		filter.State,
		"SP",
		"State",
	)

	assertStringPointerEquals(
		t,
		filter.City,
		"Garça",
		"City",
	)
}

func TestListOccurrencesEmptyFiltersBecomeNil(t *testing.T) {
	empty := "   "

	repository := &mockRepository{}

	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	_, err := service.List(
		context.Background(),
		ListOccurrencesFilter{
			Status:         &empty,
			OccurrenceType: &empty,
			AnimalType:     &empty,
			State:          &empty,
			City:           &empty,
		},
	)

	if err != nil {
		t.Fatalf(
			"não esperava erro, recebeu %v",
			err,
		)
	}

	filter := repository.listFilter

	if filter.Status != nil {
		t.Error(
			"Status deveria ser nil",
		)
	}

	if filter.OccurrenceType != nil {
		t.Error(
			"OccurrenceType deveria ser nil",
		)
	}

	if filter.AnimalType != nil {
		t.Error(
			"AnimalType deveria ser nil",
		)
	}

	if filter.State != nil {
		t.Error(
			"State deveria ser nil",
		)
	}

	if filter.City != nil {
		t.Error(
			"City deveria ser nil",
		)
	}
}

func TestListOccurrencesInvalidOccurrenceType(
	t *testing.T,
) {
	occurrenceType := "invalid"

	repository := &mockRepository{}

	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	result, err := service.List(
		context.Background(),
		ListOccurrencesFilter{
			OccurrenceType: &occurrenceType,
		},
	)

	if !errors.Is(
		err,
		ErrInvalidOccurrenceType,
	) {
		t.Errorf(
			"esperava ErrInvalidOccurrenceType, recebeu %v",
			err,
		)
	}

	if result != nil {
		t.Error(
			"resultado deveria ser nil",
		)
	}

	if repository.listCalled {
		t.Error(
			"Repository.List não deveria ser chamado",
		)
	}
}

func TestListOccurrencesInvalidState(t *testing.T) {
	state := "SPP"

	repository := &mockRepository{}

	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	result, err := service.List(
		context.Background(),
		ListOccurrencesFilter{
			State: &state,
		},
	)

	if !errors.Is(
		err,
		ErrInvalidState,
	) {
		t.Errorf(
			"esperava ErrInvalidState, recebeu %v",
			err,
		)
	}

	if result != nil {
		t.Error(
			"resultado deveria ser nil",
		)
	}

	if repository.listCalled {
		t.Error(
			"Repository.List não deveria ser chamado",
		)
	}
}

func TestListOccurrencesEmptyResult(t *testing.T) {
	repository := &mockRepository{
		listResult: []Occurrence{},
		listTotal:  0,
	}

	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	result, err := service.List(
		context.Background(),
		ListOccurrencesFilter{
			Page:     1,
			PageSize: 20,
		},
	)

	if err != nil {
		t.Fatalf(
			"não esperava erro, recebeu %v",
			err,
		)
	}

	if result == nil {
		t.Fatal(
			"esperava resultado, recebeu nil",
		)
	}

	if len(result.Items) != 0 {
		t.Errorf(
			"esperava 0 ocorrências, recebeu %d",
			len(result.Items),
		)
	}

	if result.Total != 0 {
		t.Errorf(
			"esperava total 0, recebeu %d",
			result.Total,
		)
	}

	if result.TotalPages != 0 {
		t.Errorf(
			"esperava total_pages 0, recebeu %d",
			result.TotalPages,
		)
	}
}

func TestListOccurrencesRepositoryError(t *testing.T) {
	expectedError := errors.New(
		"erro simulado ao listar ocorrências",
	)

	repository := &mockRepository{
		listErr: expectedError,
	}

	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	result, err := service.List(
		context.Background(),
		ListOccurrencesFilter{},
	)

	if !errors.Is(
		err,
		expectedError,
	) {
		t.Errorf(
			"esperava erro do repository, recebeu %v",
			err,
		)
	}

	if result != nil {
		t.Error(
			"resultado deveria ser nil quando Repository.List falha",
		)
	}

	if !repository.listCalled {
		t.Error(
			"Repository.List deveria ser chamado",
		)
	}
}
