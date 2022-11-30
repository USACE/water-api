package models

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	TimeseriesGroupFilter struct {
		IDs      *[]uuid.UUID // not supported as query param;
		Slug     *string      `param:"timeseries_group"`
		Provider *string      `param:"provider"`
	}

	TimeseriesGroup struct {
		Provider     string `json:"provider"`
		ProviderName string `json:"provider_name"`
		Slug         string `json:"slug"`
		Name         string `json:"name"`
	}

	TimeseriesGroupDetail struct {
		TimeseriesGroup
		Timeseries []struct {
			Provider string `json:"provider"`
			Datatype string `json:"datatype"`
			Key      string `json:"key"`
		} `json:"timeseries"` // abbreviated timeseries information
	}
)

func ListTimeseriesGroupsQuery(f *TimeseriesGroupFilter) (sq.SelectBuilder, error) {

	q := sq.Select(`provider, provider_name, slug, name`).From("v_timeseries_group")

	// Apply Filters
	if f != nil {
		// Filter by Slug
		if f.Slug != nil {
			q = q.Where("slug = LOWER(?)", f.Slug)
		}
		// Filter by Provider
		if f.Provider != nil {
			q = q.Where("provider = LOWER(?)", f.Provider)
		}
		// Filter by list of known UUIDs
		// This is used after CreateTimeseriesGroups when UUIDs of newly created timeseries groups are known
		if f.IDs != nil {
			q = q.Where(sq.Eq{"id": f.IDs})
		}
	}

	q = q.OrderBy("provider, slug")

	return q.PlaceholderFormat(sq.Dollar), nil
}

func ListTimeseriesGroups(db *pgxpool.Pool, f *TimeseriesGroupFilter) ([]TimeseriesGroup, error) {

	q, err := ListTimeseriesGroupsQuery(f)
	if err != nil {
		return make([]TimeseriesGroup, 0), err
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return make([]TimeseriesGroup, 0), err
	}

	vv := make([]TimeseriesGroup, 0)
	if err := pgxscan.Select(context.Background(), db, &vv, sql, args...); err != nil {
		return make([]TimeseriesGroup, 0), err
	}
	return vv, nil
}

// GetTimeseriesGroupDetail returns detailed information a single TimeseriesGroup
func GetTimeseriesGroupDetail(db *pgxpool.Pool, f *TimeseriesGroupFilter) (*TimeseriesGroupDetail, error) {

	var d TimeseriesGroupDetail
	if err := pgxscan.Get(
		context.Background(), db, &d,
		`SELECT provider, provider_name, slug, name, timeseries
		 FROM v_timeseries_group_detail
		 WHERE slug = LOWER($1)`, f.Slug,
	); err != nil {
		return nil, err
	}

	return &d, nil
}
