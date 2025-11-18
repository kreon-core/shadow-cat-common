package dbc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kreon-core/shadow-cat-common/utlc"
)

type PostgresConfig struct {
	DSN             string         `mapstructure:"dsn"                validate:"required"`
	MaxConns        *int32         `mapstructure:"max_conns"          validate:"omitempty,gte=0"`
	MinConns        *int32         `mapstructure:"min_conns"          validate:"omitempty,gte=0,ltefield=MaxConns"`
	MinIdleConns    *int32         `mapstructure:"min_idle_conns"     validate:"omitempty,gte=0,ltefield=MaxConns"`
	MaxConnIdleTime *time.Duration `mapstructure:"max_conn_idle_time" validate:"omitempty,gte=0"`
	MaxConnLifetime *time.Duration `mapstructure:"max_conn_lifetime"  validate:"omitempty,gte=0"`
}

const (
	maxConns        int32 = 10
	minConns        int32 = 1
	minIdleConns    int32 = 3
	maxConnIdleTime       = 300 * time.Second
	maxConnLifetime       = 3600 * time.Second
)

func NewPostgresConnection(ctx context.Context, cfg *PostgresConfig) (*pgxpool.Pool, error) {
	dsn := cfg.DSN
	if utlc.IsBlank(dsn) {
		return nil, errors.New(
			"DSN is blank: PostgresConfig.DSN is required but was empty",
		)
	}

	pgCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse_dsn -> %w", err)
	}

	pgCfg.MaxConns = utlc.OrElse(cfg.MaxConns, maxConns)
	pgCfg.MinConns = utlc.OrElse(cfg.MinConns, minConns)
	pgCfg.MinIdleConns = utlc.OrElse(cfg.MinIdleConns, minIdleConns)
	pgCfg.MaxConnIdleTime = utlc.OrElse(cfg.MaxConnIdleTime, maxConnIdleTime)
	pgCfg.MaxConnLifetime = utlc.OrElse(cfg.MaxConnLifetime, maxConnLifetime)

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

func ParseUUID(uuidStr string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(uuidStr)
	return uuid, err
}
