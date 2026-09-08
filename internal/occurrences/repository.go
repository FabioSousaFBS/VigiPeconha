package occurrences

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(
		ctx context.Context,
		params CreateOccurrenceParams,
	) (*Occurrence, error)
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

func (r *PostgresRepository) Create(
	ctx context.Context,
	params CreateOccurrenceParams,
) (*Occurrence, error) {
	query := `
		INSERT INTO occurrences (
			organization_id,
			created_by_user_id,
			source,
			occurrence_type,
			animal_type,
			species,
			description,
			reporter_name,
			reporter_phone,
			location,
			address,
			neighborhood,
			city,
			state,
			status,
			occurred_at
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			ST_SetSRID(
				ST_MakePoint($10, $11),
				4326
			)::geography,
			$12,
			$13,
			$14,
			$15,
			$16,
			$17
		)
		RETURNING
			id,
			organization_id,
			created_by_user_id,
			source,
			occurrence_type,
			animal_type,
			species,
			description,
			reporter_name,
			reporter_phone,
			ST_Y(location::geometry),
			ST_X(location::geometry),
			address,
			neighborhood,
			city,
			state,
			status,
			occurred_at,
			created_at,
			updated_at
	`

	var occurrence Occurrence

	err := r.db.QueryRow(
		ctx,
		query,
		params.OrganizationID,
		params.CreatedByUserID,
		params.Source,
		params.OccurrenceType,
		params.AnimalType,
		params.Species,
		params.Description,
		params.ReporterName,
		params.ReporterPhone,

		// PostGIS usa POINT(longitude latitude)
		params.Longitude,
		params.Latitude,

		params.Address,
		params.Neighborhood,
		params.City,
		params.State,
		params.Status,
		params.OccurredAt,
	).Scan(
		&occurrence.ID,
		&occurrence.OrganizationID,
		&occurrence.CreatedByUserID,
		&occurrence.Source,
		&occurrence.OccurrenceType,
		&occurrence.AnimalType,
		&occurrence.Species,
		&occurrence.Description,
		&occurrence.ReporterName,
		&occurrence.ReporterPhone,
		&occurrence.Latitude,
		&occurrence.Longitude,
		&occurrence.Address,
		&occurrence.Neighborhood,
		&occurrence.City,
		&occurrence.State,
		&occurrence.Status,
		&occurrence.OccurredAt,
		&occurrence.CreatedAt,
		&occurrence.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"erro ao criar ocorrência: %w",
			err,
		)
	}

	return &occurrence, nil
}
