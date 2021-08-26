package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PGStore struct {
	Config     Config
	Connection *pgxpool.Pool
}

func NewStore(cfg Config) (*PGStore, error) {
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
	return &PGStore{Config: cfg, Connection: db}, nil
}
