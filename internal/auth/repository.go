package auth

import (
	"context"

	"github.com/FabioSousaFBS/vigipeconha-api/internal/users"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateUser(ctx context.Context, user users.User) error
	FindByEmail(ctx context.Context, email string) (*users.User, error)
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user users.User) error {
	query := `
		INSERT INTO users (
			id,
			name,
			email,
			phone,
			password_hash,
			role
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		user.Phone,
		user.PasswordHash,
		user.Role,
	)

	return err
}

func (r *PostgresRepository) FindByEmail(ctx context.Context, email string) (*users.User, error) {
	query := `
		SELECT 
			id,
			name,
			email,
			phone,
			password_hash,
			role,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	var user users.User

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
