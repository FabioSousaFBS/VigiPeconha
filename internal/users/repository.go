package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetLicenseUsage(
		ctx context.Context,
		organizationID string,
	) (*LicenseUsage, error)

	CreateWithLicenseValidation(
		ctx context.Context,
		params CreateUserParams,
	) (*User, error)

	UpdateStatusWithLicenseValidation(
		ctx context.Context,
		userID string,
		organizationID string,
		status string,
	) (*User, error)
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

func (r *PostgresRepository) GetLicenseUsage(
	ctx context.Context,
	organizationID string,
) (*LicenseUsage, error) {
	query := `
		WITH active_contract AS (
			SELECT user_limit
			FROM contracts
			WHERE organization_id = $1
			  AND status = 'active'
			  AND payment_status = 'paid'
			  AND start_date <= CURRENT_DATE
			  AND end_date >= CURRENT_DATE
			ORDER BY end_date DESC
			LIMIT 1
		)
		SELECT
			ac.user_limit,
			COUNT(u.id) FILTER (
				WHERE u.status = 'active'
			)::INTEGER
		FROM active_contract ac
		LEFT JOIN users u
			ON u.organization_id = $1
		GROUP BY ac.user_limit
	`

	var userLimit int
	var activeUsers int

	err := r.db.QueryRow(
		ctx,
		query,
		organizationID,
	).Scan(
		&userLimit,
		&activeUsers,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrActiveContractNotFound
	}

	if err != nil {
		return nil, err
	}

	return &LicenseUsage{
		UserLimit:   userLimit,
		ActiveUsers: activeUsers,
		Available:   max(userLimit-activeUsers, 0),
	}, nil
}

func (r *PostgresRepository) CreateWithLicenseValidation(
	ctx context.Context,
	params CreateUserParams,
) (*User, error) {
	tx, err := r.db.Begin(ctx)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := lockOrganization(
		ctx,
		tx,
		params.OrganizationID,
	); err != nil {
		return nil, err
	}

	userLimit, err := getActiveContractLimit(
		ctx,
		tx,
		params.OrganizationID,
	)

	if err != nil {
		return nil, err
	}

	activeUsers, err := countActiveUsers(
		ctx,
		tx,
		params.OrganizationID,
	)

	if err != nil {
		return nil, err
	}

	if activeUsers >= userLimit {
		return nil, ErrLicenseLimitReached
	}

	query := `
		INSERT INTO users (
			name,
			email,
			phone,
			password_hash,
			role,
			status,
			organization_id
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			'active',
			$6
		)
		RETURNING
			id,
			name,
			email,
			phone,
			password_hash,
			role,
			status,
			organization_id,
			created_at,
			updated_at
	`

	var user User

	err = tx.QueryRow(
		ctx,
		query,
		params.Name,
		params.Email,
		params.Phone,
		params.PasswordHash,
		params.Role,
		params.OrganizationID,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.Role,
		&user.Status,
		&user.OrganizationID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" {
			return nil, ErrEmailAlreadyExists
		}

		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresRepository) UpdateStatusWithLicenseValidation(
	ctx context.Context,
	userID string,
	organizationID string,
	status string,
) (*User, error) {
	tx, err := r.db.Begin(ctx)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := lockOrganization(
		ctx,
		tx,
		organizationID,
	); err != nil {
		return nil, err
	}

	var currentStatus string

	err = tx.QueryRow(
		ctx,
		`
			SELECT status
			FROM users
			WHERE id = $1
			  AND organization_id = $2
			FOR UPDATE
		`,
		userID,
		organizationID,
	).Scan(&currentStatus)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	if status == "active" &&
		currentStatus != "active" {

		userLimit, err := getActiveContractLimit(
			ctx,
			tx,
			organizationID,
		)

		if err != nil {
			return nil, err
		}

		activeUsers, err := countActiveUsers(
			ctx,
			tx,
			organizationID,
		)

		if err != nil {
			return nil, err
		}

		if activeUsers >= userLimit {
			return nil, ErrLicenseLimitReached
		}
	}

	var user User

	err = tx.QueryRow(
		ctx,
		`
			UPDATE users
			SET
				status = $1,
				updated_at = NOW()
			WHERE id = $2
			  AND organization_id = $3
			RETURNING
				id,
				name,
				email,
				phone,
				password_hash,
				role,
				status,
				organization_id,
				created_at,
				updated_at
		`,
		status,
		userID,
		organizationID,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.Role,
		&user.Status,
		&user.OrganizationID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &user, nil
}

func lockOrganization(
	ctx context.Context,
	tx pgx.Tx,
	organizationID string,
) error {
	_, err := tx.Exec(
		ctx,
		`
			SELECT pg_advisory_xact_lock(
				hashtextextended($1, 0)
			)
		`,
		organizationID,
	)

	return err
}

func getActiveContractLimit(
	ctx context.Context,
	tx pgx.Tx,
	organizationID string,
) (int, error) {
	var userLimit int

	err := tx.QueryRow(
		ctx,
		`
			SELECT user_limit
			FROM contracts
			WHERE organization_id = $1
			  AND status = 'active'
			  AND payment_status = 'paid'
			  AND start_date <= CURRENT_DATE
			  AND end_date >= CURRENT_DATE
			ORDER BY end_date DESC
			LIMIT 1
			FOR UPDATE
		`,
		organizationID,
	).Scan(&userLimit)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrActiveContractNotFound
	}

	if err != nil {
		return 0, err
	}

	return userLimit, nil
}

func countActiveUsers(
	ctx context.Context,
	tx pgx.Tx,
	organizationID string,
) (int, error) {
	var count int

	err := tx.QueryRow(
		ctx,
		`
			SELECT COUNT(*)
			FROM users
			WHERE organization_id = $1
			  AND status = 'active'
		`,
		organizationID,
	).Scan(&count)

	return count, err
}
