package contracts

import (
	"context"
	"errors"

	"github.com/FabioSousaFBS/vigipeconha-api/internal/organizations"
)

var (
	ErrOrganizationInactive = errors.New(
		"organização inativa",
	)

	ErrContractUnavailable = errors.New(
		"contrato não está disponível",
	)

	ErrContractOverdue = errors.New(
		"contrato com pagamento em atraso",
	)
)

type Service interface {
	ValidateOrganizationAccess(
		ctx context.Context,
		organizationID string,
	) error
}

type ContractService struct {
	contractRepository     Repository
	organizationRepository organizations.Repository
}

func NewService(
	contractRepository Repository,
	organizationRepository organizations.Repository,
) Service {
	return &ContractService{
		contractRepository:     contractRepository,
		organizationRepository: organizationRepository,
	}
}

func (s *ContractService) ValidateOrganizationAccess(
	ctx context.Context,
	organizationID string,
) error {
	organization, err :=
		s.organizationRepository.FindStatusByID(
			ctx,
			organizationID,
		)

	if err != nil {
		return err
	}

	if organization.Status != "active" {
		return ErrOrganizationInactive
	}

	contract, err :=
		s.contractRepository.FindActiveByOrganizationID(
			ctx,
			organizationID,
		)

	if err != nil {
		if errors.Is(err, ErrNoActiveContract) {
			return ErrContractUnavailable
		}

		return err
	}

	if contract.PaymentStatus == "overdue" {
		return ErrContractOverdue
	}

	if contract.PaymentStatus != "paid" {
		return ErrContractUnavailable
	}

	return nil
}
