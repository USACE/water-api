package models

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Provider struct {
	Slug string `json:"provider" db:"slug"`
	Name string `json:"provider_name" db:"name"`
}

func ListProviders(db *pgxpool.Pool) ([]Provider, error) {
	pp := make([]Provider, 0)
	if err := pgxscan.Select(
		context.Background(), db, &pp, `SELECT slug, name FROM provider ORDER BY slug`,
	); err != nil {
		return make([]Provider, 0), err
	}
	return pp, nil
}
