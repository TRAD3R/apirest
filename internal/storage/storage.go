package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trad3r/hskills/apirest/internal/domain"
	"github.com/trad3r/hskills/apirest/internal/repository/postgres"
)

type Storage struct {
	User domain.UserRepository
	Post domain.PostRepository
}

var conn *pgxpool.Pool

func NewDB(ctx context.Context, url string) (*Storage, error) {
	conn, err := pgxpool.New(ctx, url)
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
