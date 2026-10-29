package healthfacilities

import "errors"

var (
	ErrInvalidInput  = errors.New("dados da unidade inválidos")
	ErrInvalidID     = errors.New("identificador da unidade inválido")
	ErrNotFound      = errors.New("unidade de atendimento não encontrada")
	ErrDuplicateCNES = errors.New("CNES já cadastrado")
)
