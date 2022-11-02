package models

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Datasource struct {
	Provider     string `json:"provider"`
	ProviderName string `json:"provider_name"`
	Datatype     string `json:"datatype"`
	DatatypeName string `json:"datatype_name"`
	DatatypeUri  string `json:"datatype_uri"`
}

// ListDatasources
func ListDatasources(db *pgxpool.Pool, p string, d string) ([]Datasource, error) {
	q := sq.Select(
		`p.slug    AS provider,
		 p.name    AS provider_name,
		 dt.slug   AS datatype,
		 dt.name   AS datatype_name,
		 dt.uri    AS datatype_uri`,
	).From(
		"datasource AS d",
	).Join(
		"provider AS p ON p.id = d.provider_id",
	).Join(
		"datatype AS dt ON dt.id = d.datatype_id",
	).PlaceholderFormat(sq.Dollar)

	if p != "" {
		q = q.Where("p.Slug = ?", p)
	}
	if d != "" {
		q = q.Where("dt.Slug = ?", d)
	}

	sql, args, err := q.ToSql()

	dd := make([]Datasource, 0)
	if err != nil {
		return dd, err
	}
	if err := pgxscan.Select(
		context.Background(), db, &dd, sql, args...); err != nil {
		return dd, err
	}
	return dd, nil
}
