package occurrences

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(
		ctx context.Context,
		params CreateOccurrenceParams,
	) (*Occurrence, error)

	CreatePhoto(
		ctx context.Context,
		photo OccurrencePhoto,
	) (*OccurrencePhoto, error)

	Exists(
		ctx context.Context,
		occurrenceID string,
	) (bool, error)

	List(
		ctx context.Context,
		filter ListOccurrencesFilter,
	) ([]Occurrence, int64, error)
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

func (r *PostgresRepository) Exists(
	ctx context.Context,
	occurrenceID string,
) (bool, error) {
	var exists bool

	err := r.db.QueryRow(
		ctx,
		`
		SELECT EXISTS (
			SELECT 1
			FROM occurrences
			WHERE id = $1
		)
		`,
		occurrenceID,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf(
			"erro ao verificar ocorrência: %w",
			err,
		)
	}

	return exists, nil
}

func (r *PostgresRepository) CreatePhoto(
	ctx context.Context,
	photo OccurrencePhoto,
) (*OccurrencePhoto, error) {
	query := `
		INSERT INTO occurrence_photos (
			occurrence_id,
			storage_key,
			photo_url,
			content_type,
			file_size
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
		)
		RETURNING
			id,
			occurrence_id,
			storage_key,
			photo_url,
			content_type,
			file_size,
			created_at
	`

	var created OccurrencePhoto

	err := r.db.QueryRow(
		ctx,
		query,
		photo.OccurrenceID,
		photo.StorageKey,
		photo.PhotoURL,
		photo.ContentType,
		photo.FileSize,
	).Scan(
		&created.ID,
		&created.OccurrenceID,
		&created.StorageKey,
		&created.PhotoURL,
		&created.ContentType,
		&created.FileSize,
		&created.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"erro ao salvar foto da ocorrência: %w",
			err,
		)
	}

	return &created, nil
}

func (r *PostgresRepository) List(
	ctx context.Context,
	filter ListOccurrencesFilter,
) ([]Occurrence, int64, error) {
	whereClause, args := buildListWhereClause(filter)

	countQuery := `
		SELECT COUNT(*)
		FROM occurrences
	` + whereClause

	var total int64

	if err := r.db.QueryRow(
		ctx,
		countQuery,
		args...,
	).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf(
			"erro ao contar ocorrências: %w",
			err,
		)
	}

	offset := (filter.Page - 1) * filter.PageSize

	queryArgs := append(
		[]any{},
		args...,
	)

	limitPosition := len(queryArgs) + 1
	queryArgs = append(queryArgs, filter.PageSize)

	offsetPosition := len(queryArgs) + 1
	queryArgs = append(queryArgs, offset)

	query := fmt.Sprintf(`
		SELECT
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
		FROM occurrences
		%s
		ORDER BY occurred_at DESC, created_at DESC
		LIMIT $%d
		OFFSET $%d
	`,
		whereClause,
		limitPosition,
		offsetPosition,
	)

	rows, err := r.db.Query(
		ctx,
		query,
		queryArgs...,
	)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"erro ao listar ocorrências: %w",
			err,
		)
	}
	defer rows.Close()

	occurrences := make([]Occurrence, 0)

	for rows.Next() {
		var occurrence Occurrence

		err := rows.Scan(
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
			return nil, 0, fmt.Errorf(
				"erro ao ler ocorrência: %w",
				err,
			)
		}

		occurrences = append(
			occurrences,
			occurrence,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf(
			"erro ao percorrer ocorrências: %w",
			err,
		)
	}

	return occurrences, total, nil
}

func buildListWhereClause(
	filter ListOccurrencesFilter,
) (string, []any) {
	conditions := make([]string, 0)
	args := make([]any, 0)

	addCondition := func(
		column string,
		value *string,
		caseInsensitive bool,
	) {
		if value == nil {
			return
		}

		args = append(args, *value)

		position := len(args)

		if caseInsensitive {
			conditions = append(
				conditions,
				fmt.Sprintf(
					"LOWER(%s) = LOWER($%d)",
					column,
					position,
				),
			)

			return
		}

		conditions = append(
			conditions,
			fmt.Sprintf(
				"%s = $%d",
				column,
				position,
			),
		)
	}

	addCondition(
		"status",
		filter.Status,
		false,
	)

	addCondition(
		"occurrence_type",
		filter.OccurrenceType,
		false,
	)

	addCondition(
		"animal_type",
		filter.AnimalType,
		true,
	)

	addCondition(
		"state",
		filter.State,
		true,
	)

	addCondition(
		"city",
		filter.City,
		true,
	)

	if len(conditions) == 0 {
		return "", args
	}

	return " WHERE " +
			strings.Join(conditions, " AND "),
		args
}
