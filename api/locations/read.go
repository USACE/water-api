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
	Slug     *string      `param:"location"` // if set, should always return 0 or 1 locations
	State    *string      `query:"state"`
	Provider *string      `query:"provider"`
	Datatype *string      `query:"datatype"`
	Code     *string      `query:"code"`
	Q        *string      `query:"q"`
}

func ListLocationsQuery(f *LocationFilter) (sq.SelectBuilder, error) {

	q := sq.Select(
		`provider,
		 provider_name,
		 datatype,
		 datatype_name,
		 code,
		 slug,
		 ST_AsGeoJSON(geometry)::json AS geometry,
		 state,
		 state_name,
		 attributes`,
	).From(
		"v_location l",
	)

	// Apply Filters (excluding Search Query String)
	if f != nil {
		// Filter by Slug
		if f.Slug != nil {
			q = q.Where("slug = lower(?)", f.Slug)
		}
		// Filter by State
		if f.State != nil {
			q = q.Where("state = UPPER(?)", f.State)
		}

		// Filter by Provider
		if f.Provider != nil {
			q = q.Where("provider = LOWER(?)", f.Provider)
		}

		// Filter by Datatype
		if f.Datatype != nil {
			// WHERE
			q = q.Where("datatype = LOWER(?)", f.Datatype)
		}

		// Filter by Code
		if f.Code != nil {
			q = q.Where("code = LOWER(?)", f.Code)
		}

		// Filter by Search String
		if f.Q != nil {
			q = q.Where("slug || code ILIKE '%' || ? || '%'", f.Q)
		}

		// Filter by list of known UUIDs
		// This is used after CreateLocations when UUIDs of newly created locations are known
		if f.IDs != nil {
			q = q.Where(sq.Eq{"id": f.IDs})
		}
	}

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

func GetLocation(db *pgxpool.Pool, f *LocationFilter) (*LocationInfo, error) {

	q, err := ListLocationsQuery(f)
	if err != nil {
		return nil, err
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	var l LocationInfo
	if err := pgxscan.Get(context.Background(), db, &l, sql, args...); err != nil {
		return nil, err
	}

	return &l, nil
}
