package cwms

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	Config     Config
	Connection *pgxpool.Pool
}

type Config struct {
	ApplicationKey        string        `envconfig:"APPLICATION_KEY"`
	DBUser                string        `envconfig:"DB_USER"`
	DBPass                string        `envconfig:"DB_PASS"`
	DBName                string        `envconfig:"DB_NAME"`
	DBHost                string        `envconfig:"DB_HOST"`
	DBSSLMode             string        `envconfig:"DB_SSL_MODE"`
	DBPoolMaxConns        int           `envconfig:"DB_POOL_MAX_CONNS" default:"10"`
	DBPoolMaxConnIdleTime time.Duration `envconfig:"DB_POOL_MAX_CONN_IDLE_TIME" default:"30m"`
	DBPoolMinConns        int           `envconfig:"DB_POOL_MIN_CONNS" default:"5"`
}

func NewStore(cfg Config) (*Store, error) {
	poolConfig, err := pgxpool.ParseConfig(
		fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s sslmode=%s",
			cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBHost, cfg.DBSSLMode,
		),
	)
	if err != nil {
		return nil, err
	}
	poolConfig.MaxConns = int32(cfg.DBPoolMaxConns)
	poolConfig.MinConns = int32(cfg.DBPoolMinConns)
	poolConfig.MaxConnIdleTime = cfg.DBPoolMaxConnIdleTime

	db, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}
	return &Store{Config: cfg, Connection: db}, nil
}
