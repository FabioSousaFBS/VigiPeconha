package users

import "errors"

var (
	ErrUserNotFound = errors.New(
		"usuário não encontrado",
	)

	ErrEmailAlreadyExists = errors.New(
		"já existe um usuário cadastrado com este e-mail",
	)

	ErrLicenseLimitReached = errors.New(
		"limite de usuários ativos do contrato atingido",
	)

	ErrActiveContractNotFound = errors.New(
		"contrato ativo não encontrado",
	)

	ErrInvalidRole = errors.New(
		"perfil de usuário inválido",
	)

	ErrInvalidStatus = errors.New(
		"status de usuário inválido",
	)

	ErrInvalidPassword = errors.New(
		"a senha deve possuir pelo menos 8 caracteres",
	)

	ErrInvalidName = errors.New(
		"nome é obrigatório",
	)

	ErrInvalidEmail = errors.New(
		"e-mail é obrigatório",
	)
)
