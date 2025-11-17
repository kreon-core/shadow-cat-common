package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	tul "github.com/kreon-core/shadow-cat-common"
)

const (
	maxConns        int32 = 10
	minConns        int32 = 1
	minIdleConns    int32 = 3
	maxConnIdleTime       = 300 * time.Second
	maxConnLifetime       = 3600 * time.Second
)

func NewConnection(ctx context.Context, cfg *Config) (*pgxpool.Pool, error) {
	dsn := cfg.DSN
	if tul.IsBlank(dsn) {
		return nil, errors.New("empty dsn string")
	}

	pgCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse_dsn -> %w", err)
	}

	pgCfg.MaxConns = tul.OrElse(cfg.MaxConns, maxConns)
	pgCfg.MinConns = tul.OrElse(cfg.MinConns, minConns)
	pgCfg.MinIdleConns = tul.OrElse(cfg.MinIdleConns, minIdleConns)
	pgCfg.MaxConnIdleTime = tul.OrElse(cfg.MaxConnIdleTime, maxConnIdleTime)
	pgCfg.MaxConnLifetime = tul.OrElse(cfg.MaxConnLifetime, maxConnLifetime)

	pool, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		return nil, fmt.Errorf("create_connection_pool -> %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping_database -> %w", err)
	}

	return pool, nil
}
