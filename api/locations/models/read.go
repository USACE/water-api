package models

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

type LocationFilter struct {
	Slugs    *[]string // not supported as query param at this time
	Slug     *string   `query:"location" param:"location"` // binds to either /locations/:slug or /locations?slug=
	State    *string   `query:"state"`
	Provider *string   `query:"provider"`
	Q        *string   `query:"q"`
}

func ListLocationsQuery(f *LocationFilter) (sq.SelectBuilder, error) {

	q := sq.Select(
		`p.slug                         AS provider_slug,
	     p.name                         AS provider_name,
		 a.slug                         AS slug,
	     a.name                         AS name,
		 s.abbreviation                 AS state,
		 ST_AsGeoJSON(a.geometry)::json AS geometry`,
	).From(
		"location_v2 a",
	)

	// Basic Join Statements
	jS := "tiger.state s ON s.gid         = a.state_id"    // join tiger.state
	jD := "datasource  d ON d.provider_id = a.provider_id" // join datasource
	jP := "provider    p ON p.id          = d.provider_id" // join provider

	// Apply Filters (excluding Search Query String)
	if f != nil {
		// Filter by State
		if f.State != nil {
			q = q.Join(jS+" AND s.abbreviation = ?", f.State)
			q = q.Where("s.abbreviation = ?", f.State)
		} else {
			q = q.Join(jS) // always join tiger.state
		}

		// Filter by Provider
		if f.Provider != nil {
			q = q.Join(jD) // Join datasource. todo; additional qualifier needed
			q = q.Join(jP) // Join provider. todo; additioinal qualifier needed
			q = q.Where("p.slug = ?", f.Provider)
		} else {
			q = q.Join(jD)
			q = q.Join(jP)
		}

		// Filter by Search String
		if f.Q != nil {
			q = q.Where("a.slug || a.name ILIKE '%' || ? || '%'", f.Q) // filter by query string
		}
	} else {
		q = q.Join(jS).Join(jD).Join(jP) // always join state, datasource, provider tables even if no filters in url
	}

	return q.PlaceholderFormat(sq.Dollar), nil
}

func ListLocations(db *pgxpool.Pool, f *LocationFilter) ([]Location, error) {
	q, err := ListLocationsQuery(f)
	if err != nil {
		return make([]Location, 0), err
	}
	sql, args, err := q.ToSql()
	if err != nil {
		return make([]Location, 0), err
	}
	ll := make([]Location, 0)
	if err := pgxscan.Select(context.Background(), db, &ll, sql, args...); err != nil {
		return make([]Location, 0), err
	}
	return ll, nil
}

func GetLocation(db *pgxpool.Pool, f *LocationFilter) (*Location, error) {
	q, err := ListLocationsQuery(f)
	if err != nil {
		return nil, err
	}
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	var l Location
	if err := pgxscan.Get(context.Background(), db, &l, sql, args...); err != nil {
		return nil, err
	}
	return &l, nil
}
