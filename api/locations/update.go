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
			 SET geometry = $1,
			     state_id = (SELECT gid FROM tiger_data.state_all WHERE stusps = UPPER($2)),
			     attributes = $3
			 WHERE datasource_id = (
				SELECT id
				FROM datasource
				WHERE datatype_id = (SELECT id FROM datatype WHERE slug = LOWER($4)) AND provider_id = (SELECT id FROM provider WHERE slug = LOWER($5))
			 )
			 AND code = $6
			 RETURNING id`,
			info.Geometry.EWKT(6), info.State, info.Attributes, info.Datatype, info.Provider, info.Code,
		)
		if err != nil {
			tx.Rollback(context.Background())
			return make([]LocationInfo, 0), err
		}
		var id uuid.UUID
		if err := pgxscan.ScanOne(&id, rows); err != nil {
			tx.Rollback(context.Background())
			return make([]LocationInfo, 0), err
		} // todo; test coverage and confirm behavior when UPDATE called on site that does not exist.

		IDs = append(IDs, id)
	}
	tx.Commit(context.Background())

	return ListLocations(db, &LocationFilter{IDs: &IDs})
}
