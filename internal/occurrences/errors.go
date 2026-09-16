package occurrences

import "errors"

var (
	ErrInvalidOccurrenceType = errors.New(
		"tipo de ocorrência inválido",
	)

	ErrAnimalTypeRequired = errors.New(
		"tipo do animal é obrigatório",
	)

	ErrCoordinatesRequired = errors.New(
		"latitude e longitude são obrigatórias",
	)

	ErrInvalidLatitude = errors.New(
		"latitude inválida",
	)

	ErrInvalidLongitude = errors.New(
		"longitude inválida",
	)

	ErrOccurredAtRequired = errors.New(
		"data da ocorrência é obrigatória",
	)

	ErrInvalidState = errors.New(
		"estado deve possuir 2 caracteres",
	)

	ErrOccurrenceNotFound = errors.New(
		"ocorrência não encontrada",
	)

	ErrEmptyPhoto       = errors.New("arquivo vazio")
	ErrPhotoTooLarge    = errors.New("imagem deve possuir no máximo 5 MB")
	ErrInvalidPhotoType = errors.New("formato de imagem não permitido")
	ErrInvalidPhoto     = errors.New("não foi possível processar a imagem")
)
