package models

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

func DeleteLocation(db *pgxpool.Pool, locationID *uuid.UUID) error {
	if _, err := db.Exec(context.Background(), `DELETE FROM location WHERE id = $1`, locationID); err != nil {
		return err
	}
	return nil
}
