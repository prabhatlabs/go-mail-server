package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prabhatlabs/go-mail-server/internal/db"
)

type DB struct {
	Pool *pgxpool.Pool
	Q    *db.Queries
}

func ConnectDB(ctx context.Context, url string) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{
		Pool: pool,
		Q: db.New(pool),
	}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}

func (db *DB) WithTx(ctx context.Context, f func(tx *db.Queries) error) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := f(db.Q.WithTx(tx)); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
