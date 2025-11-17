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
	maxConns        = 10
	minConns        = 1
	minIdleConns    = 3
	maxConnIdleTime = 300 * time.Second
	maxConnLifetime = 3600 * time.Second
)

func CreateConnection(ctx context.Context, cfg *Config) (*pgxpool.Pool, error) {
	dsn := cfg.DSN
	if tul.IsBlank(dsn) {
		return nil, errors.New("empty dsn string")
	}

	pgCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse_dsn -> %w", err)
	}

	pgCfg.MaxConns = int32(tul.OrElse(cfg.MaxConns, int(maxConns)))
	pgCfg.MinConns = int32(tul.OrElse(cfg.MinConns, int(minConns)))
	pgCfg.MinIdleConns = int32(tul.OrElse(cfg.MinIdleConns, int(minIdleConns)))
	pgCfg.MaxConnIdleTime = tul.OrElse(cfg.MaxConnIdleTime, maxConnIdleTime)
	pgCfg.MaxConnLifetime = tul.OrElse(cfg.MaxConnLifetime, maxConnLifetime)

	pool, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		return nil, fmt.Errorf("create_connection_pool -> %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping_database -> %w", err)
	}

	return pool, nil
}
