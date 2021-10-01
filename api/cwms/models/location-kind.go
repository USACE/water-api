package models

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type LocationKind struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// type LocationKindCollection struct {
// 	Items []LocationKind `json:"items`
// }

func ListLocationKind(db *pgxpool.Pool) ([]LocationKind, error) {
	lk := make([]LocationKind, 0)
	if err := pgxscan.Select(
		context.Background(),
		db,
		&lk,
		"SELECT k.id, k.name FROM location_kind k",
	); err != nil {
		return make([]LocationKind, 0), err
	}
	return lk, nil
}
