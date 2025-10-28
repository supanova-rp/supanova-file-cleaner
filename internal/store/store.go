package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/supanova-rp/supanova-file-cleaner/internal/store/sqlc"
)

type Store struct {
	pool    *pgxpool.Pool
	Queries *sqlc.Queries
}

func NewStore(ctx context.Context, dbUrl string) (*Store, error) {
	poolConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %v", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %v", err)
	}

	return &Store{
		pool:    pool,
		Queries: sqlc.New(pool),
	}, nil
}

func (s *Store) Close() {
	s.pool.Close()
}
