package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	pgx "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trad3r/hskills/apirest/internal/domain"
	"github.com/trad3r/hskills/apirest/internal/repository/postgres"
	"log"
)

type Storage struct {
	User domain.UserRepository
	Post domain.PostRepository
}

var conn *pgxpool.Pool

func NewDB(ctx context.Context, dsn string) (*Storage, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse config: %v", err)
	}

	// Настройка параметров пула соединений
	config.MaxConns = 15
	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("could not ping database: %w", err)
	}

	if err := runMigration(dsn); err != nil {
		return nil, err
	}

	return &Storage{User: postgres.NewUserRepository(conn), Post: postgres.NewPostRepository(conn)}, nil
}

func runMigration(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}

	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return fmt.Errorf("could not connect get database driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("could not init migration: %w", err)
	}
	if err := m.Up(); err != nil {
		return fmt.Errorf("could not run migration: %w", err)
	}

	return nil
}

func Stop() {
	conn.Close()
}
