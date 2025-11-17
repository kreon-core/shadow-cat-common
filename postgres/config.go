package postgres

import "time"

type Config struct {
	DSN             string         `mapstructure:"dsn"                validate:"required"`
	MaxConns        *int32         `mapstructure:"max_conns"          validate:"omitempty,gte=0"`
	MinConns        *int32         `mapstructure:"min_conns"          validate:"omitempty,gte=0,ltefield=MaxConns"`
	MinIdleConns    *int32         `mapstructure:"min_idle_conns"     validate:"omitempty,gte=0,ltefield=MaxConns"`
	MaxConnIdleTime *time.Duration `mapstructure:"max_conn_idle_time" validate:"omitempty,gte=0"`
	MaxConnLifetime *time.Duration `mapstructure:"max_conn_lifetime"  validate:"omitempty,gte=0"`
}
