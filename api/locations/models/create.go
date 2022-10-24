package models

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Example POST Body Payload (cwms_location)
//
//	{
//	  "slug": "location-1",
//	  "geometry": {geojson},
//	  "state": "mn",
//	  "provider": "mvp",
//	  "type": "cwms-location"
//	}
//
// Example POST Body Payload (nws_site)
//
//	{
//	  "slug": "nws-location-1",
//	  "geometry": {geojson},
//	  "state": "mn",
//	  "provider": "nws-serfc",
//	  "type": "nws-site"
//	}
//
// Example POST Body Payload (usgs_site)
//
//	{
//	  "slug": "01090150",    // unique string
//	  "geometry": {geojson},
//	  "state": "mn",         // fk back to state_id
//	  "provider": "usgs",    // fk back to provider; along with type, determines datasource
//	  "type": "usgs-site"    // fk back to datasource_type; along with provider, determines datasource
//	}
func Create(db *pgxpool.Pool, n LocationCollection) ([]Location, error) {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]Location, 0), err
	}
	defer tx.Rollback(context.Background())
	newIDs := make([]uuid.UUID, 0)
	for _, m := range n.Items {
		rows, err := tx.Query(
			context.Background(),
			`INSERT INTO location_v2 (datasource_id, name, slug, geometry) VALUES ($1, $2, $3, $4, $5, $6) RETURNING slug`,
			m.OfficeID, m.Name, m.PublicName, m.Slug, m.Geometry.EWKT(6), m.KindID,
		)
		if err != nil {
			tx.Rollback(context.Background())
			return make([]Location, 0), err
		}
		var id uuid.UUID
		if err := pgxscan.ScanOne(&id, rows); err != nil {
			tx.Rollback(context.Background())
			return make([]Location, 0), err
		} else {
			newIDs = append(newIDs, id)
		}
	}
	tx.Commit(context.Background())

	return ListLocationsForIDs(db, newIDs)
}
