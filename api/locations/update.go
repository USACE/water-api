package locations

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

func (cc LocationCollection) Update(db *pgxpool.Pool) ([]LocationInfo, error) {

	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]LocationInfo, 0), err
	}
	defer tx.Rollback(context.Background())

	IDs := make([]uuid.UUID, 0)
	for _, n := range cc.Items {
		info := n.LocationInfo()
		rows, err := tx.Query(
			context.Background(),
			`UPDATE location
			 SET geometry = ST_GeomFromGeoJSON($1::json),
			     state_id = (SELECT gid FROM tiger_data.state_all WHERE stusps = UPPER($2)),
			     attributes = $3
			 WHERE datasource_id = (
				SELECT id
				FROM datasource
				WHERE datatype_id = (SELECT id FROM datatype WHERE slug = LOWER($4)) AND provider_id = (SELECT id FROM provider WHERE slug = LOWER($5))
			 )
			 AND code = LOWER($6)
			 RETURNING id`,
			info.Geometry, info.State, info.Attributes, info.Datatype, info.Provider, info.Code,
		)
		if err != nil {
			tx.Rollback(context.Background())
			return make([]LocationInfo, 0), err
		}
		var id uuid.UUID
		if err := pgxscan.ScanOne(&id, rows); err != nil {
			// Row was not updated; Do not include ID in list of updated records
			continue
		}
		IDs = append(IDs, id)
	}
	tx.Commit(context.Background())

	return ListLocations(db, &LocationFilter{IDs: &IDs})
}
