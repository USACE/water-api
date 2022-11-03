package locations

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type LocationFilter struct {
	IDs      *[]uuid.UUID // not supported as query param;
	State    *string      `query:"state"`
	Provider *string      `query:"provider"`
	Datatype *string      `query:"datatype"`
	Code     *string      `query:"code"`
	Q        *string      `query:"q"`
}

func ListLocationsQuery(f *LocationFilter) (sq.SelectBuilder, error) {

	q := sq.Select(
		`p.slug                         AS provider,
	     p.name                         AS provider_name,
		 dt.slug                        AS datatype,
		 dt.name                        AS datatype_name,
		 l.code                         AS code,
		 l.slug                         AS slug,
		 ST_AsGeoJSON(l.geometry)::json AS geometry,
		 s.stusps                       AS state,
		 l.attributes                   AS attributes`,
	).From(
		"location l",
	)

	// Join Statements
	jS, jSParams := "tiger.state  s  ON s.gid  = l.state_id", make([]interface{}, 0)      // join tiger.state
	jDS, jDSParams := "datasource ds ON ds.id  = l.datasource_id", make([]interface{}, 0) // join datasource
	jP, jPParams := "provider     p  ON p.id   = ds.provider_id", make([]interface{}, 0)  // join provider
	jDT, jDTParams := "datatype   dt ON dt.id  = ds.datatype_id", make([]interface{}, 0)  // join datatype

	// Apply Filters (excluding Search Query String)
	if f != nil {
		// Filter by State
		if f.State != nil {
			// Limit JOIN using ?state= query param
			jS += " AND s.gid = (SELECT gid FROM tiger.state WHERE stusps = UPPER(?))"
			jSParams = append(jSParams, f.State)
			// WHERE
			q = q.Where("s.gid = (SELECT gid FROM tiger.state WHERE stusps = UPPER(?))", f.State)
		}

		// Filter by Provider
		if f.Provider != nil {
			// Limit datasource join by ?provider= query param
			jDS += " AND ds.provider_id = (SELECT id from provider WHERE slug = LOWER(?))"
			jDSParams = append(jDSParams, f.Provider)
			// Limit provider join by ?provider= query param
			jP += " AND p.slug = LOWER(?)"
			jPParams = append(jPParams, f.Provider)
			// WHERE
			q = q.Where("p.slug = ?", f.Provider)
		}

		// Filter by Datatype
		if f.Datatype != nil {
			// Limit datasource join by ?datatype= query param
			jDS += " AND ds.datatype_id = (SELECT id from datatype WHERE slug = LOWER(?))"
			jDSParams = append(jDSParams, f.Datatype)
			// Limit datatype join by ?datatype= query param
			jDT += " AND dt.slug = LOWER(?)"
			jDTParams = append(jDTParams, f.Datatype)
			// WHERE
			q = q.Where("dt.slug = ?", f.Datatype)
		}

		// Filter by Code
		if f.Code != nil {
			q = q.Where("UPPER(l.code) = UPPER(?)", f.Code) // todo; confirm case insensitive is desired behavior
		}

		// Filter by Search String
		if f.Q != nil {
			q = q.Where("l.slug || l.code ILIKE '%' || ? || '%'", f.Q) // filter by query string
		}

		// Filter by list of known UUIDs
		// This is used after CreateLocations when UUIDs of newly created locations are known
		if f.IDs != nil {
			q = q.Where(sq.Eq{"l.id": f.IDs})
		}
	}

	q = q.Join(jS, jSParams...)   // join state
	q = q.Join(jDS, jDSParams...) // join datasource
	q = q.Join(jP, jPParams...)   // join provider
	q = q.Join(jDT, jDTParams...) // join datatype

	return q.PlaceholderFormat(sq.Dollar), nil
}

func ListLocations(db *pgxpool.Pool, f *LocationFilter) ([]LocationInfo, error) {
	q, err := ListLocationsQuery(f)
	if err != nil {
		return make([]LocationInfo, 0), err
	}
	sql, args, err := q.ToSql()
	if err != nil {
		return make([]LocationInfo, 0), err
	}
	ll := make([]LocationInfo, 0)
	if err := pgxscan.Select(context.Background(), db, &ll, sql, args...); err != nil {
		return make([]LocationInfo, 0), err
	}
	return ll, nil
}
