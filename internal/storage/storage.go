package storage

import (
	"context"
	"fmt"
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trad3r/hskills/apirest/internal/domain"
	"github.com/trad3r/hskills/apirest/internal/repository/postgres"
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

	return &Storage{User: postgres.NewUserRepository(conn), Post: postgres.NewPostRepository(conn)}, nil
}

func Stop() {
	conn.Close()
}
