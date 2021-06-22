package models

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

type State struct {
	Abbreviation string `json:"abbreviation"`
	Name         string `json:"name"`
}

func ListStates(db *pgxpool.Pool) ([]State, error) {
	ss := make([]State, 0)
	if err := pgxscan.Select(
		context.Background(), db, &ss,
		`SELECT stusps AS abbreviation, name
		 FROM tiger.state
		 ORDER BY name`,
	); err != nil {
		return make([]State, 0), err
	}
	return ss, nil
}
