package locations

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Delete(db *pgxpool.Pool, ic LocationInfoCollection) error {
	for _, l := range ic.Items {
		if _, err := db.Exec(
			context.Background(),
			`DELETE FROM location
			 WHERE (
				datasource_id = (SELECT id FROM v_datasource WHERE datatype = LOWER($1) AND provider = LOWER($2))
				AND (LOWER(code) = LOWER($3) OR slug = LOWER($4))
			 ) OR (
				slug = LOWER($4) AND datasource_id IN (SELECT id FROM v_datasource WHERE provider = LOWER($2))
			 )`, l.Datatype, l.Provider, l.Code, l.Slug,
		); err != nil {
			return err
		}
	}
	return nil
}
