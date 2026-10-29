package occurrences

import (
	"context"
	"errors"
	"testing"
)

const testOccurrenceID = "550e8400-e29b-41d4-a716-446655440000"

func TestGetByIDInvalidUUID(t *testing.T) {
	repository := &mockRepository{}
	service := NewService(repository, &mockStorage{})

	result, err := service.GetByID(
		context.Background(),
		"abc",
	)

	if !errors.Is(err, ErrInvalidOccurrenceID) {
		t.Fatalf("esperava ErrInvalidOccurrenceID, recebeu %v", err)
	}

	if result != nil {
		t.Fatal("resultado deveria ser nil")
	}

	if repository.findByIDCalled || repository.listPhotosCalled {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestGetByIDOccurrenceNotFound(t *testing.T) {
	repository := &mockRepository{
		findByIDErr: ErrOccurrenceNotFound,
	}
	service := NewService(repository, &mockStorage{})

	result, err := service.GetByID(
		context.Background(),
		testOccurrenceID,
	)

	if !errors.Is(err, ErrOccurrenceNotFound) {
		t.Fatalf("esperava ErrOccurrenceNotFound, recebeu %v", err)
	}

	if result != nil {
		t.Fatal("resultado deveria ser nil")
	}

	if !repository.findByIDCalled {
		t.Fatal("FindByID deveria ser chamado")
	}

	if repository.listPhotosCalled {
		t.Fatal("ListPhotos não deveria ser chamado")
	}
}

func TestGetByIDWithPhotos(t *testing.T) {
	repository := &mockRepository{
		findByIDResult: &Occurrence{
			ID:         testOccurrenceID,
			AnimalType: "snake",
			Status:     "pending",
		},
		listPhotosResult: []OccurrencePhoto{
			{
				ID:           "photo-1",
				OccurrenceID: testOccurrenceID,
				StorageKey:   "occurrences/photo-1.jpg",
				ContentType:  "image/jpeg",
			},
			{
				ID:           "photo-2",
				OccurrenceID: testOccurrenceID,
				StorageKey:   "occurrences/photo-2.png",
				ContentType:  "image/png",
			},
		},
	}

	service := NewService(repository, &mockStorage{})

	result, err := service.GetByID(
		context.Background(),
		testOccurrenceID,
	)

	if err != nil {
		t.Fatalf("não esperava erro: %v", err)
	}

	if result == nil {
		t.Fatal("resultado não deveria ser nil")
	}

	if result.ID != testOccurrenceID {
		t.Errorf("ID inesperado: %s", result.ID)
	}

	if len(result.Photos) != 2 {
		t.Fatalf("esperava 2 fotos, recebeu %d", len(result.Photos))
	}

	if result.Photos[0].ID != "photo-1" {
		t.Errorf("primeira foto incorreta: %s", result.Photos[0].ID)
	}

	if !repository.findByIDCalled || !repository.listPhotosCalled {
		t.Fatal("FindByID e ListPhotos deveriam ser chamados")
	}
}

func TestGetByIDWithoutPhotos(t *testing.T) {
	repository := &mockRepository{
		findByIDResult: &Occurrence{
			ID: testOccurrenceID,
		},
		listPhotosResult: []OccurrencePhoto{},
	}

	service := NewService(repository, &mockStorage{})

	result, err := service.GetByID(
		context.Background(),
		testOccurrenceID,
	)

	if err != nil {
		t.Fatalf("não esperava erro: %v", err)
	}

	if result == nil {
		t.Fatal("resultado não deveria ser nil")
	}

	if result.Photos == nil {
		t.Fatal("Photos deveria ser uma lista vazia, não nil")
	}

	if len(result.Photos) != 0 {
		t.Fatalf("esperava 0 fotos, recebeu %d", len(result.Photos))
	}
}

func TestGetByIDFindByIDRepositoryError(t *testing.T) {
	expectedError := errors.New("erro ao consultar ocorrência")

	repository := &mockRepository{
		findByIDErr: expectedError,
	}

	service := NewService(repository, &mockStorage{})

	result, err := service.GetByID(
		context.Background(),
		testOccurrenceID,
	)

	if !errors.Is(err, expectedError) {
		t.Fatalf("esperava erro do repository, recebeu %v", err)
	}

	if result != nil {
		t.Fatal("resultado deveria ser nil")
	}

	if repository.listPhotosCalled {
		t.Fatal("ListPhotos não deveria ser chamado")
	}
}

func TestGetByIDListPhotosRepositoryError(t *testing.T) {
	expectedError := errors.New("erro ao consultar fotos")

	repository := &mockRepository{
		findByIDResult: &Occurrence{
			ID: testOccurrenceID,
		},
		listPhotosErr: expectedError,
	}

	service := NewService(repository, &mockStorage{})

	result, err := service.GetByID(
		context.Background(),
		testOccurrenceID,
	)

	if !errors.Is(err, expectedError) {
		t.Fatalf("esperava erro de ListPhotos, recebeu %v", err)
	}

	if result != nil {
		t.Fatal("resultado deveria ser nil")
	}

	if !repository.findByIDCalled || !repository.listPhotosCalled {
		t.Fatal("ambos os métodos deveriam ser chamados")
	}
}
