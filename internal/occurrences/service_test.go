package occurrences

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
)

type mockRepository struct {
	exists bool

	createPhotoCalled bool
	createPhotoErr    error
}

func (m *mockRepository) Create(
	ctx context.Context,
	params CreateOccurrenceParams,
) (*Occurrence, error) {
	return nil, nil
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
	// Arrange
	repository := &mockRepository{}
	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	const fileSize int64 = 5*1024*1024 + 1

	// Act
	_, err := service.UploadPhoto(
		context.Background(),
		"occurrence-id",
		fileSize,
		bytes.NewReader(nil),
	)

	// Assert
	if !errors.Is(err, ErrPhotoTooLarge) {
		t.Errorf(
			"esperava ErrPhotoTooLarge, recebeu %v",
			err,
		)
	}
}

func TestUploadPhotoInvalidContentType(t *testing.T) {
	// Arrange
	repository := &mockRepository{
		exists: true,
	}
	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	data := []byte("este arquivo nao e uma imagem")

	// Act
	_, err := service.UploadPhoto(
		context.Background(),
		"occurrence-id",
		int64(len(data)),
		bytes.NewReader(data),
	)

	// Assert
	if !errors.Is(err, ErrInvalidPhotoType) {
		t.Errorf(
			"esperava ErrInvalidPhotoType, recebeu %v",
			err,
		)
	}
}

func TestUploadPhotoOccurrenceNotFound(t *testing.T) {
	// Arrange
	repository := &mockRepository{
		exists: false,
	}

	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	pngData := validPNGData()

	// Act
	_, err := service.UploadPhoto(
		context.Background(),
		"occurrence-id",
		int64(len(pngData)),
		bytes.NewReader(pngData),
	)

	// Assert
	if !errors.Is(err, ErrOccurrenceNotFound) {
		t.Errorf(
			"esperava ErrOccurrenceNotFound, recebeu %v",
			err,
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

func TestUploadPhotoSuccess(t *testing.T) {
	// Arrange
	repository := &mockRepository{
		exists: true,
	}

	storage := &mockStorage{}

	service := NewService(
		repository,
		storage,
	)

	pngData := validPNGData()

	// Act
	photo, err := service.UploadPhoto(
		context.Background(),
		"occurrence-id",
		int64(len(pngData)),
		bytes.NewReader(pngData),
	)

	// Assert
	if err != nil {
		t.Fatalf(
			"não esperava erro, recebeu %v",
			err,
		)
	}

	if photo == nil {
		t.Fatal("esperava uma foto, recebeu nil")
	}

	if !storage.uploadCalled {
		t.Error("esperava que Storage.Upload fosse chamado")
	}

	if !repository.createPhotoCalled {
		t.Error("esperava que Repository.CreatePhoto fosse chamado")
	}

	if storage.deleteCalled {
		t.Error("Storage.Delete não deveria ser chamado")
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
}

func TestUploadPhotoStorageError(t *testing.T) {
	// Arrange
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

	// Act
	_, err := service.UploadPhoto(
		context.Background(),
		"occurrence-id",
		int64(len(pngData)),
		bytes.NewReader(pngData),
	)

	// Assert
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
	// Arrange
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

	// Act
	_, err := service.UploadPhoto(
		context.Background(),
		"occurrence-id",
		int64(len(pngData)),
		bytes.NewReader(pngData),
	)

	// Assert
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
