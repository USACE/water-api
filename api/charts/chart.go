package charts

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
	ChartFilter struct {
		IDs      *[]uuid.UUID // not supported as query param;
		Slug     *string      `param:"chart"`
		Provider *string      `query:"provider" param:"provider"`
	}

	Chart struct {
		Provider     string `json:"provider"`
		ProviderName string `json:"provider_name" db:"provider_name"`
		Type         string `json:"type"`
		Slug         string `json:"slug"`
		Name         string `json:"name"`
		Location     *struct {
			Slug     *string `json:"slug"`     // Optional Location Information; required to establish linkage to unique location on Create
			Provider *string `json:"provider"` // Optional Location Information; required to establish linkage to unique location on Create
			Datatype *string `json:"datatype"` // Optional Location Information; required to establish linkage to unique location on Create
			Code     *string `json:"code"`     // Optional Location Information; required to establish linkage to unique location on Create
		} `json:"location"`
		Mapping *[]ChartMapping `json:"mapping,omitempty"`
	}

	ChartCollection struct {
		Items []Chart `json:"items"`
	}
)

func (c *ChartCollection) UnmarshalJSON(b []byte) error {
	switch helpers.JSONType(b) {
	case "ARRAY":
		return json.Unmarshal(b, &c.Items)
	case "OBJECT":
		c.Items = make([]Chart, 1)
		return json.Unmarshal(b, &c.Items[0])
	default:
		return errors.New("payload not recognized as JSON array or object")
	}
}

func ListChartsQuery(f *ChartFilter) (sq.SelectBuilder, error) {

	q := sq.Select(`location, slug, name, type, provider, provider_name`).From("v_chart")

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
		// This is used after CreateLocations when UUIDs of newly created locations are known
		if f.IDs != nil {
			q = q.Where(sq.Eq{"id": f.IDs})
		}
	}

	q = q.OrderBy("provider, slug")

	return q.PlaceholderFormat(sq.Dollar), nil
}

func ListCharts(db *pgxpool.Pool, f *ChartFilter) ([]Chart, error) {

	q, err := ListChartsQuery(f)
	if err != nil {
		return make([]Chart, 0), err
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return make([]Chart, 0), err
	}

	vv := make([]Chart, 0)
	if err := pgxscan.Select(context.Background(), db, &vv, sql, args...); err != nil {
		return make([]Chart, 0), err
	}
	return vv, nil
}

// GetChartDetail returns detailed information a single Chart
func GetChartDetail(db *pgxpool.Pool, f *ChartFilter) (*Chart, error) {

	var t Chart
	if err := pgxscan.Get(
		context.Background(), db, &t,
		`SELECT location, slug, name, type, provider, provider_name, mapping
		 FROM v_chart_detail
		 WHERE slug = LOWER($1)`, f.Slug,
	); err != nil {
		return nil, err
	}

	return &t, nil
}

// CreateChart creates a single Chart
func (cc ChartCollection) Create(db *pgxpool.Pool) ([]Chart, error) {

	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]Chart, 0), err
	}
	defer tx.Rollback(context.Background())
	newIDs := make([]uuid.UUID, 0)

	// Create a map of all existing slugs in the database.
	// Append the map each time a new location is created and another slug is taken.
	slugMap, err := helpers.SlugMap(db, "chart", "slug", "", "")
	if err != nil {
		return make([]Chart, 0), err
	}
	for _, n := range cc.Items {
		// Get Unique Slug for Each Location
		slug, err := helpers.GetUniqueSlug(n.Name, slugMap)
		if err != nil {
			return make([]Chart, 0), err
		}
		// Add slug to map so it's not reused within this transaction
		slugMap[slug] = true

		rows, err := tx.Query(
			context.Background(),
			`INSERT INTO chart (slug, name, type, provider_id, location_id)
			 VALUES
			 	(
					$1,
					$2,
					LOWER($3),
					(SELECT id from provider WHERE slug = LOWER($4)),
					(SELECT id FROM location WHERE slug = LOWER($5))
				)
			RETURNING id`, slug, n.Name, n.Type, n.Provider, n.Location.Slug,
		)
		if err != nil {
			tx.Rollback(context.Background())
			return make([]Chart, 0), err
		}
		var id uuid.UUID
		if err := pgxscan.ScanOne(&id, rows); err != nil {
			continue // Chart already exists; DO NOTHING bypasses RETURNING id
		}
		newIDs = append(newIDs, id)
	}
	tx.Commit(context.Background())

	return ListCharts(db, &ChartFilter{IDs: &newIDs})
}

func DeleteChart(db *pgxpool.Pool, chartProvider *string, chartSlug *string) error {
	if _, err := db.Exec(
		context.Background(),
		`DELETE FROM chart 
		 WHERE provider_id = (SELECT id FROM provider WHERE slug = LOWER($1))
		   AND slug = LOWER($2)`, chartProvider, chartSlug,
	); err != nil {
		return err
	}
	return nil
}
