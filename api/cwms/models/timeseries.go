package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/USACE/water-api/api/helpers"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Timeseries struct {
	ID             *uuid.UUID  `json:"id,omitempty"`
	Provider       string      `json:"provider" db:"provider"`
	DatasourceType string      `json:"datasource_type" db:"datasource_type"`
	Key            string      `json:"key"`
	Times          []time.Time `json:"times" db:"times"`
	Values         []float64   `json:"values" db:"values"`
}

type TimeseriesFilter struct {
	DatasourceType *string `json:"datasource_type" query:"datasource_type"`
	Provider       *string `query:"provider"`
	Q              *string `query:"q"`
}

type TimeseriesCollection struct {
	Items []Timeseries `json:"items"`
}

func (c *TimeseriesCollection) UnmarshalJSON(b []byte) error {
	switch helpers.JSONType(b) {
	case "ARRAY":
		return json.Unmarshal(b, &c.Items)
	case "OBJECT":
		c.Items = make([]Timeseries, 1)
		return json.Unmarshal(b, &c.Items[0])
	default:
		return errors.New("payload not recognized as JSON array or object")
	}
}

func ListTimeseriesQuery(f *TimeseriesFilter) (sq.SelectBuilder, error) {

	q := sq.Select(`dt.slug AS datasource_type,
					p.slug AS provider,
					t.datasource_key AS key,
					ARRAY_AGG(t.latest_time) AS times,
					ARRAY_AGG(t.latest_value ) AS values`,
	).From("timeseries t")

	// Base string for JOIN
	j1 := `datasource d ON d.id = t.datasource_id 
			JOIN datasource_type dt ON dt.id = d.datasource_type_id 
			JOIN provider p ON p.id = d.provider_id`

	q = q.Join(j1)

	q = q.GroupBy("dt.slug, p.slug, t.datasource_key")

	if f != nil {

		// Filter by Provider
		if f.Provider != nil {
			q = q.Where("lower(p.slug) = lower(?)", *f.Provider)
		}
		// Filter by DatasourceType
		if f.DatasourceType != nil {
			q = q.Where("lower(dt.slug) = lower(?)", *f.DatasourceType)
		}
	}

	fmt.Println(q.ToSql())

	// Unfiltered
	return q.PlaceholderFormat(sq.Dollar), nil
}

func ListTimeseries(db *pgxpool.Pool, f *TimeseriesFilter) ([]Timeseries, error) {
	q, err := ListTimeseriesQuery(f)
	if err != nil {
		return make([]Timeseries, 0), err
	}
	sql, args, err := q.ToSql()
	if err != nil {
		return make([]Timeseries, 0), err
	}
	tt := make([]Timeseries, 0)
	if err := pgxscan.Select(context.Background(), db, &tt, sql, args...); err != nil {
		return make([]Timeseries, 0), err
	}
	return tt, nil
}

func CreateOrUpdateTimeseries(db *pgxpool.Pool, c TimeseriesCollection) ([]Timeseries, error) {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]Timeseries, 0), err
	}
	defer tx.Rollback(context.Background())

	newIDs := make([]uuid.UUID, 0)

	queryDataSourceID := `
		SELECT d.id FROM datasource d 
		JOIN datasource_type dt ON dt.id = d.datasource_type_id 
		JOIN provider p ON p.id = d.provider_id 
		WHERE lower(p.slug) = lower($2) AND lower(dt.slug) = lower($1)
	`

	for _, t := range c.Items {
		rows, err := tx.Query(
			context.Background(),
			`INSERT INTO timeseries (datasource_id, datasource_key, latest_time, latest_value)
			VALUES((`+queryDataSourceID+`), $3, $4, $5)
			ON CONFLICT ON CONSTRAINT timeseries_unique_datasource
			DO UPDATE SET
			latest_time = $4,
			latest_value = $5
			WHERE timeseries.datasource_key = $3
			AND timeseries.datasource_id = (`+queryDataSourceID+`)
			RETURNING id`,
			t.DatasourceType, t.Provider, t.Key, t.Times[len(t.Times)-1], t.Values[len(t.Values)-1],
		)
		if err != nil {
			return make([]Timeseries, 0), err
		}
		var id uuid.UUID
		if err := pgxscan.ScanOne(&id, rows); err != nil {
			tx.Rollback(context.Background())
			return c.Items, err
		} else {
			newIDs = append(newIDs, id)
		}
	}
	tx.Commit(context.Background())
	return nil, nil
}
