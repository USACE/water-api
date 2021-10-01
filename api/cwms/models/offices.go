package models

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Office struct {
	ID       uuid.UUID  `json:"id"`
	Name     string     `json:"name"`
	Symbol   string     `json:"symbol"`
	ParentID *uuid.UUID `json:"parent_id" db:"parent_id"`
}

func ListOffices(db *pgxpool.Pool) ([]Office, error) {
	oo := make([]Office, 0)
	if err := pgxscan.Select(
		context.Background(), db, &oo,
		`SELECT id, name, symbol, parent_id
		 FROM office
		 ORDER BY symbol`,
	); err != nil {
		return make([]Office, 0), err
	}
	return oo, nil
}
