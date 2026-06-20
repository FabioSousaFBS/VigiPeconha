package database

import (
	"context"
	"fmt"
	"log"

	"github.com/FabioSousaFBS/vigipeconha-api/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, cfg config.Config) *pgxpool.Pool {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("erro ao criar pool do banco: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("erro ao conectar no banco: %v", err)
	}

	log.Println("conexão com PostgreSQL realizada com sucesso")

	return pool
}
