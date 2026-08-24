package contracts

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNoActiveContract = errors.New("contrato ativo não encontrado")

type ActiveContract struct {
	ID             string
	OrganizationID string
	StartDate      time.Time
	EndDate        time.Time
	UserLimit      int
	Status         string
	PaymentStatus  string
}

type Repository interface {
	FindActiveByOrganizationID(
		ctx context.Context,
		organizationID string,
	) (*ActiveContract, error)
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) Repository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) FindActiveByOrganizationID(
	ctx context.Context,
	organizationID string,
) (*ActiveContract, error) {
	query := `
		SELECT
			id,
			organization_id,
			start_date,
			end_date,
			user_limit,
			status,
			payment_status
		FROM contracts
		WHERE organization_id = $1
		  AND status = 'active'
		  AND start_date <= CURRENT_DATE
		  AND end_date >= CURRENT_DATE
		ORDER BY end_date DESC
		LIMIT 1
	`

	var contract ActiveContract

	err := r.db.QueryRow(
		ctx,
		query,
		organizationID,
	).Scan(
		&contract.ID,
		&contract.OrganizationID,
		&contract.StartDate,
		&contract.EndDate,
		&contract.UserLimit,
		&contract.Status,
		&contract.PaymentStatus,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNoActiveContract
	}

	if err != nil {
		return nil, err
	}

	return &contract, nil
}
