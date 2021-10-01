package models

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Parameter struct {
	UID         uuid.UUID `json:"-" db:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
}

func ListParameters(db *pgxpool.Pool) ([]Parameter, error) {
	pp := make([]Parameter, 0)
	if err := pgxscan.Select(
		context.Background(), db, &pp,
		`SELECT id, code, description
		 FROM usgs_parameter`,
	); err != nil {
		return make([]Parameter, 0), err
	}
	return pp, nil
}
