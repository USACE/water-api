package locations

import (
	"context"

	"github.com/USACE/water-api/api/helpers"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

func (cc LocationCollection) Create(db *pgxpool.Pool) ([]LocationInfo, error) {

	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]LocationInfo, 0), err
	}
	defer tx.Rollback(context.Background())
	newIDs := make([]uuid.UUID, 0)

	// Create a map of all existing slugs in the database.
	// Append the map each time a new location is created and another slug is taken.
	slugMap, err := helpers.SlugMap(db, "location", "slug", "", "")
	if err != nil {
		return make([]LocationInfo, 0), err
	}
	for _, n := range cc.Items {
		info := n.LocationInfo()
		// Get Unique Slug for Each Location
		slug, err := helpers.GetUniqueSlug(info.Code, slugMap)
		if err != nil {
			return make([]LocationInfo, 0), err
		}
		// Add slug to map so it's not reused within this transaction
		slugMap[slug] = true

		rows, err := tx.Query(
			context.Background(),
			`INSERT INTO location (datasource_id, slug, code, geometry, state_id, attributes)
			 VALUES (
			    (
					SELECT id 
				      FROM v_datasource
				     WHERE datatype = LOWER($1) AND provider = LOWER($2)
				),
				$3,
				$4,
				ST_GeomFromGeoJSON($5::json),
				(
					SELECT gid
				      FROM tiger_data.state_all
				     WHERE stusps = UPPER($6)
				),
				$7
			 ) ON CONFLICT ON CONSTRAINT datasource_unique_code DO NOTHING
			 RETURNING id`,
			info.Datatype, info.Provider, slug, info.Code, info.Geometry, info.State, info.Attributes,
		)
		if err != nil {
			tx.Rollback(context.Background())
			return make([]LocationInfo, 0), err
		}
		var id uuid.UUID
		if err := pgxscan.ScanOne(&id, rows); err != nil {
			continue // Location already exists for given datasource and code; DO NOTHING on constraint bypasses RETURNING id
		}
		newIDs = append(newIDs, id)
	}
	tx.Commit(context.Background())

	return ListLocations(db, &LocationFilter{IDs: &newIDs})
}
