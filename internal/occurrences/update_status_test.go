package occurrences

import (
	"context"
	"errors"
	"testing"
)

const (
	statusTestOccurrenceID = "550e8400-e29b-41d4-a716-446655440000"
	statusTestUserID       = "550e8400-e29b-41d4-a716-446655440001"
)

func TestUpdateStatusSuccess(t *testing.T) {
	repository := &mockRepository{}
	service := NewService(repository, &mockStorage{})

	result, err := service.UpdateStatus(
		context.Background(),
		statusTestOccurrenceID,
		"  ADDRESSED  ",
		statusTestUserID,
	)

	if err != nil {
		t.Fatalf("não esperava erro: %v", err)
	}

	if !repository.updateStatusCalled {
		t.Fatal("UpdateStatus deveria ser chamado")
	}

	if repository.updateStatusValue != StatusAddressed {
		t.Fatalf(
			"esperava addressed, recebeu %s",
			repository.updateStatusValue,
		)
	}

	if repository.updateStatusUserID != statusTestUserID {
		t.Fatal("usuário responsável incorreto")
	}

	if result.Status != StatusAddressed {
		t.Fatalf("status retornado incorreto: %s", result.Status)
	}
}

func TestUpdateStatusInvalidUUID(t *testing.T) {
	repository := &mockRepository{}
	service := NewService(repository, &mockStorage{})

	_, err := service.UpdateStatus(
		context.Background(),
		"abc",
		StatusAddressed,
		statusTestUserID,
	)

	if !errors.Is(err, ErrInvalidOccurrenceID) {
		t.Fatalf("erro inesperado: %v", err)
	}

	if repository.updateStatusCalled {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestUpdateStatusInvalidStatus(t *testing.T) {
	repository := &mockRepository{}
	service := NewService(repository, &mockStorage{})

	_, err := service.UpdateStatus(
		context.Background(),
		statusTestOccurrenceID,
		"resolved",
		statusTestUserID,
	)

	if !errors.Is(err, ErrInvalidOccurrenceStatus) {
		t.Fatalf("erro inesperado: %v", err)
	}

	if repository.updateStatusCalled {
		t.Fatal("repository não deveria ser chamado")
	}
}

func TestUpdateStatusOccurrenceNotFound(t *testing.T) {
	repository := &mockRepository{
		updateStatusErr: ErrOccurrenceNotFound,
	}

	service := NewService(repository, &mockStorage{})

	_, err := service.UpdateStatus(
		context.Background(),
		statusTestOccurrenceID,
		StatusArchived,
		statusTestUserID,
	)

	if !errors.Is(err, ErrOccurrenceNotFound) {
		t.Fatalf("erro inesperado: %v", err)
	}
}

func TestUpdateStatusUnchanged(t *testing.T) {
	repository := &mockRepository{
		updateStatusErr: ErrOccurrenceStatusUnchanged,
	}

	service := NewService(repository, &mockStorage{})

	_, err := service.UpdateStatus(
		context.Background(),
		statusTestOccurrenceID,
		StatusPending,
		statusTestUserID,
	)

	if !errors.Is(err, ErrOccurrenceStatusUnchanged) {
		t.Fatalf("erro inesperado: %v", err)
	}
}

func TestUpdateStatusRepositoryError(t *testing.T) {
	expectedError := errors.New("falha no banco")

	repository := &mockRepository{
		updateStatusErr: expectedError,
	}

	service := NewService(repository, &mockStorage{})

	_, err := service.UpdateStatus(
		context.Background(),
		statusTestOccurrenceID,
		StatusAcknowledged,
		statusTestUserID,
	)

	if !errors.Is(err, expectedError) {
		t.Fatalf("erro inesperado: %v", err)
	}
}
