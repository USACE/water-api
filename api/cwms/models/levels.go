package models

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Level struct
type Level struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Slug string    `json:"slug" db:"slug"`
	Name string    `json:"name" db:"name"`
}

// ListLevelKind
func ListLevelKind(db *pgxpool.Pool) ([]Level, error) {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	tx.Rollback(context.Background())

	lvl := []Level{}
	if err = pgxscan.Select(context.Background(),
		db,
		&lvl,
		`SELECT id, slug, name FROM level_kind`,
	); err != nil {
		return nil, err
	}

	return lvl, nil
}

// CreateLocationLevel
// func CreateLocationLevel(db *pgxpool.Pool) error {
// 	return nil
// }
