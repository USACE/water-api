package locations

import (
	"context"

	"github.com/USACE/water-api/api/helpers"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	// LocationCreator defines the behaviors that must be supported for all
	// location types to create records in the database using the Create method
	LocationCreator interface {
		LocationInfo() Location
		CreateAttributes(tx *pgx.Tx, locationID *uuid.UUID) error
	}

	// LocationCreatorCollection holds LocationCreator interfaces
	// to support different behaviors for locations that have different
	// underlying datasource_type properties
	LocationCreatorCollection struct {
		Items []LocationCreator
	}
)

func (cc LocationCreatorCollection) Create(db *pgxpool.Pool) ([]Location, error) {

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
		slug, err := helpers.GetUniqueSlug(info.Slug, slugMap)
		if err != nil {
			return make([]Location, 0), err
		}
		// Add slug to map so it's not reused within this transaction
		slugMap[slug] = true

		rows, err := tx.Query(
			context.Background(),
			`INSERT INTO location (datasource_id, slug, geometry, state_id)
			 VALUES (
			    (
					SELECT id 
				      FROM datasource
				     WHERE datasource_type_id = (SELECT id FROM datasource_type WHERE slug = $1)
					   AND provider_id        = (SELECT id FROM provider WHERE slug = $2)
				),
				$3,
				$4,
				(
					SELECT gid
				      FROM tiger_data.state_all
				     WHERE UPPER(stusps) = UPPER($5)
				)
			 ) RETURNING id`,
			info.DatasourceType, info.Provider, slug, info.Geometry.EWKT(6), info.State,
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

		// CreateAttributes. Delegated to each different kind of location
		// Pass a transaction pointer and the location_id
		// and let the Createable take care of the rest.
		if err := n.CreateAttributes(&tx, &id); err != nil {
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
