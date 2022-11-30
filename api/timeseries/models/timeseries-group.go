package models

import (
	"context"
	"encoding/json"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/USACE/water-api/api/helpers"
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

	TimeseriesGroupCollection struct {
		Items []TimeseriesGroup `json:"items"`
	}

	TimeseriesGroupDetail struct {
		TimeseriesGroup
		Timeseries []TimeseriesGroupMember `json:"timeseries"` // abbreviated timeseries information
	}
)

func (c *TimeseriesGroupCollection) UnmarshalJSON(b []byte) error {
	switch helpers.JSONType(b) {
	case "ARRAY":
		return json.Unmarshal(b, &c.Items)
	case "OBJECT":
		c.Items = make([]TimeseriesGroup, 1)
		return json.Unmarshal(b, &c.Items[0])
	default:
		return errors.New("payload not recognized as JSON array or object")
	}
}

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

// CreateTimeseriesGroups creates a Timeseries Group
func CreateTimeseriesGroups(db *pgxpool.Pool, gc *TimeseriesGroupCollection) ([]TimeseriesGroup, error) {

	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]TimeseriesGroup, 0), err
	}
	defer tx.Rollback(context.Background())
	newIDs := make([]uuid.UUID, 0)

	// Create a map of all existing slugs in the database.
	// Append the map each time a new location is created and another slug is taken.
	slugMap, err := helpers.SlugMap(db, "timeseries_group", "slug", "", "")
	if err != nil {
		return make([]TimeseriesGroup, 0), err
	}
	for _, n := range gc.Items {
		// Get Unique Slug for Each Location
		slug, err := helpers.GetUniqueSlug(n.Name, slugMap)
		if err != nil {
			return make([]TimeseriesGroup, 0), err
		}
		// Add slug to map so it's not reused within this transaction
		slugMap[slug] = true

		rows, err := tx.Query(
			context.Background(),
			`INSERT INTO timeseries_group (slug, name, provider_id)
			 VALUES
			 	(
					$1,
					$2,
					(SELECT id from provider WHERE slug = LOWER($3))
				)
			 ON CONFLICT ON CONSTRAINT provider_unique_timeseries_group_name DO NOTHING
			 RETURNING id`, slug, n.Name, n.Provider,
		)
		if err != nil {
			tx.Rollback(context.Background())
			return make([]TimeseriesGroup, 0), err
		}
		var id uuid.UUID
		if err := pgxscan.ScanOne(&id, rows); err != nil {
			continue // TimeseriesGroup already exists; DO NOTHING bypasses RETURNING id
		}
		newIDs = append(newIDs, id)
	}
	tx.Commit(context.Background())

	return ListTimeseriesGroups(db, &TimeseriesGroupFilter{IDs: &newIDs})
}

func DeleteTimeseriesGroup(db *pgxpool.Pool, provider *string, slug *string) error {
	if _, err := db.Exec(
		context.Background(),
		`DELETE FROM timeseries_group 
		 WHERE provider_id = (SELECT id FROM provider WHERE slug = LOWER($1))
		   AND slug = LOWER($2)`, provider, slug,
	); err != nil {
		return err
	}
	return nil
}
