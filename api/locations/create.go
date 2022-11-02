package locations

import (
	"context"

	"github.com/USACE/water-api/api/helpers"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

func (cc LocationCollection) Create(db *pgxpool.Pool) ([]Location, error) {

	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]Location, 0), err
	}
	defer tx.Rollback(context.Background())
	newIDs := make([]uuid.UUID, 0)

	// Create a map of all existing slugs in the database.
	// Append the map each time a new location is created and another slug is taken.
	slugMap, err := helpers.SlugMap(db, "location", "slug", "", "")
	if err != nil {
		return make([]Location, 0), err
	}
	for _, n := range cc.Items {
		info := n.LocationInfo()
		// Get Unique Slug for Each Location
		slug, err := helpers.GetUniqueSlug(info.Code, slugMap)
		if err != nil {
			return make([]Location, 0), err
		}
		// Add slug to map so it's not reused within this transaction
		slugMap[slug] = true

		rows, err := tx.Query(
			context.Background(),
			`INSERT INTO location (datasource_id, slug, code, geometry, state_id, attributes)
			 VALUES (
			    (
					SELECT id 
				      FROM datasource
				     WHERE datatype_id = (SELECT id FROM datatype WHERE slug = $1)
					   AND provider_id = (SELECT id FROM provider WHERE slug = $2)
				),
				$3,
				$4,
				$5,
				(
					SELECT gid
				      FROM tiger_data.state_all
				     WHERE UPPER(stusps) = UPPER($6)
				),
				$7
			 ) RETURNING id`,
			info.Datatype, info.Provider, slug, info.Code, info.Geometry.EWKT(6), info.State, info.Attributes,
		)
		if err != nil {
			tx.Rollback(context.Background())
			return make([]Location, 0), err
		}
		var id uuid.UUID
		if err := pgxscan.ScanOne(&id, rows); err != nil {
			tx.Rollback(context.Background())
			return make([]Location, 0), err
		}

		newIDs = append(newIDs, id)
	}
	tx.Commit(context.Background())

	return make([]Location, 0), nil
	// todo
	// return ListLocationsForIDs(db, newIDs)
}
