package models

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/USACE/water-api/api/helpers"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Location struct {
	Slug         string           `json:"slug"`
	Name         string           `json:"name"`
	Geometry     helpers.Geometry `json:"geometry"`
	State        *string          `json:"state"` // state abbreviation (e.g. MN, TN, WV, FL)
	ProviderSlug string           `json:"provider_slug"`
	ProviderName string           `json:"provider_name"`
}

type LocationFilter struct {
	Slugs    *[]string // not supported as query param at this time
	Slug     *string   `query:"location" param:"location"` // binds to either /locations/:slug or /locations?slug=
	State    *string   `query:"state"`
	Provider *string   `query:"provider"`
	Q        *string   `query:"q"`
}

type LocationCollection struct {
	Items []Location `json:"items"`
}

func (c *LocationCollection) UnmarshalJSON(b []byte) error {
	switch helpers.JSONType(b) {
	case "ARRAY":
		return json.Unmarshal(b, &c.Items)
	case "OBJECT":
		c.Items = make([]Location, 1)
		return json.Unmarshal(b, &c.Items[0])
	default:
		return errors.New("payload not recognized as JSON array or object")
	}
}

// Sync Locations
func SyncLocations(db *pgxpool.Pool, c LocationCollection) ([]Location, error) {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]Location, 0), err
	}
	defer tx.Rollback(context.Background())

	newSlugs := make([]string, 0)

	for _, l := range c.Items {
		rows, err := tx.Query(
			context.Background(),
			`INSERT INTO location (datasource_id, name, slug, geometry, kind_id)
			VALUES($1, $2, $3, $4, $5, $6)
			ON CONFLICT ON CONSTRAINT office_unique_name
			DO UPDATE SET
			public_name = $3,
			geometry = $5,
			kind_id = $6,
			update_date = CURRENT_TIMESTAMP
			WHERE location.office_id = $1 AND location.name = $2
			RETURNING slug`,
			l.OfficeID, l.Name, l.PublicName, l.Slug, l.Geometry.EWKT(6), l.KindID,
		)
		if err != nil {
			return make([]Location, 0), err
		}
		var slug string
		if err := pgxscan.ScanOne(&slug, rows); err != nil {
			tx.Rollback(context.Background())
			return c.Items, err
		} else {
			newSlugs = append(newSlugs, slug)
		}
	}
	tx.Commit(context.Background())
	return ListLocations(db, &LocationFilter{Slugs: &newSlugs})
}
