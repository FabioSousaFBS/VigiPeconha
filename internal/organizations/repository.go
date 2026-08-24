package organizations

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrOrganizationNotFound = errors.New(
	"organização não encontrada",
)

type OrganizationStatus struct {
	ID     string
	Status string
}

type Repository interface {
	FindStatusByID(
		ctx context.Context,
		id string,
	) (*OrganizationStatus, error)
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(
	db *pgxpool.Pool,
) Repository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) FindStatusByID(
	ctx context.Context,
	id string,
) (*OrganizationStatus, error) {
	query := `
		SELECT
			id,
			status
		FROM organizations
		WHERE id = $1
		LIMIT 1
	`

	var organization OrganizationStatus

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&organization.ID,
		&organization.Status,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrOrganizationNotFound
	}

	if err != nil {
		return nil, err
	}

	return &organization, nil
}
