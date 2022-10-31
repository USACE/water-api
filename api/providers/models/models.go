package models

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Provider struct {
	Slug string `json:"provider" db:"slug"`
	Name string `json:"provider_name" db:"name"`
}

type Datasource struct {
	Slug string `json:"datasource_type" db:"datasource_type"`
	URI  string `json:"uri" db:"uri"`
}

type DatasourceProvider struct {
	Provider
	Datasource
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

// ListDatasources
func ListDatasources(db *pgxpool.Pool, p string, d string) ([]DatasourceProvider, error) {
	q := sq.Select(
		`p.slug         AS        provider,
		p.name          AS        provider_name,
		dt.slug         AS        datasource_type, dt.uri`,
	).From(
		"datasource AS d",
	).Join(
		"provider AS p ON p.id = d.provider_id",
	).Join(
		"datasource_type AS dt ON dt.id = d.datasource_type_id",
	).PlaceholderFormat(sq.Dollar)

	if p != "" {
		q = q.Where("p.Slug = ?", p)
	}
	if d != "" {
		q = q.Where("dt.Slug = ?", d)
	}

	sql, args, err := q.ToSql()

	dd := make([]DatasourceProvider, 0)
	if err != nil {
		return dd, err
	}
	if err := pgxscan.Select(
		context.Background(), db, &dd, sql, args...); err != nil {
		return dd, err
	}
	return dd, nil
}
