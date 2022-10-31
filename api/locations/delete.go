package locations

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

func DeleteLocation(db *pgxpool.Pool, locationID *uuid.UUID) error {
	// if _, err := db.Exec(context.Background(), `UPDATE location SET deleted=true WHERE UPPER(slug) = $1`, locationID); err != nil {
	// 	return err
	// }
	return nil
}
