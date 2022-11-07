package models

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/USACE/water-api/api/helpers"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

// POST /timeseries
//
//	Creates a Timeseries may or may not have a Measurements struct.  A measurements struct has two fields. Times, Values. Each is an array with zero or more values []
//	Creates metadata that will be associated with all related timeseries measurements
type Timeseries struct {
	ID           *uuid.UUID    `json:"id,omitempty"`
	Provider     string        `json:"provider" db:"provider"`
	Datatype     string        `json:"datatype" db:"datatype"`
	Key          string        `json:"key"`
	LocationSlug string        `json:"location_slug,omitempty"`
	LatestTime   *time.Time    `json:"latest_time,omitempty"`
	LatestValue  *float64      `json:"latest_value,omitempty"`
	Measurements *Measurements `json:"measurements,omitempty"`
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

type TimeseriesFilter struct {
	Datatype   *string `json:"datatype" query:"datatype"`
	Provider   *string `query:"provider"`
	OnlyMapped bool    `query:"only_mapped"`
	Q          *string `query:"q"`
}

func ListTimeseriesQuery(f *TimeseriesFilter) (sq.SelectBuilder, error) {

	q := sq.Select(
		`dt.slug 			AS datatype,
		 p.slug  			AS provider,
		 l.slug				AS location_slug,
		 t.datasource_key 	AS key,
		 t.latest_time,
		 t.latest_value`,
	).From(
		"timeseries t",
	).Join(
		"datasource d ON d.id = t.datasource_id",
	).Join(
		"datatype dt ON dt.id = d.datatype_id",
	).Join(
		"provider p ON p.id = d.provider_id",
	)

	q = q.Join("location l on l.id = t.location_id")

	if f != nil {

		// Filter by Provider
		if f.Provider != nil {
			q = q.Where("lower(p.slug) = lower(?)", *f.Provider)
		}
		// Filter by Datatype
		if f.Datatype != nil {
			q = q.Where("lower(dt.slug) = lower(?)", *f.Datatype)
		}
		// Filter by IsMapped
		if f.OnlyMapped {
			q = q.Join("chart_variable_mapping cvm ON cvm.timeseries_id = t.id")
		}
		// Filter by search string

		if f.Q != nil {
			if len(*f.Q) > 2 {
				q = q.Where("lower(t.datasource_key) ILIKE '%' || lower(?) || '%' ", f.Q)
			}
		}

	}

	// fmt.Println(q.ToSql())

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
	if err := pgxscan.Select(context.Background(), db, &tt, sql+" order by p.slug, t.datasource_key", args...); err != nil {
		return make([]Timeseries, 0), err
	}
	return tt, nil
}

// func GetTimeseries() (*Timeseries, error) {}

// func CreateTimeseries() ([]Timeseries, error) {}
func CreateTimeseries(db *pgxpool.Pool, c TimeseriesCollection) ([]Timeseries, error) {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]Timeseries, 0), err
	}
	defer tx.Rollback(context.Background())

	//newIDs := make([]uuid.UUID, 0)

	queryDataSourceID := `
		SELECT d.id FROM datasource d 
		JOIN datatype dt ON dt.id = d.datatype_id 
		JOIN provider p ON p.id = d.provider_id 
		WHERE lower(p.slug) = lower($2) AND lower(dt.slug) = lower($1)
	`

	for _, t := range c.Items {
		rows, err := tx.Query(
			context.Background(),
			`INSERT INTO timeseries (datasource_id, datasource_key, location_id)
			VALUES((`+queryDataSourceID+`), $3, (SELECT id from location where lower(slug) = lower($4)))	
			ON CONFLICT DO NOTHING		
			RETURNING id`,
			t.Datatype, t.Provider, t.Key, t.LocationSlug,
		)
		if err != nil {
			return make([]Timeseries, 0), err
		}
		var id uuid.UUID
		if err := pgxscan.ScanOne(&id, rows); err != nil {
			// tx.Rollback(context.Background())
			// return c.Items, err
			continue
		}
		// } else {
		// 	newIDs = append(newIDs, id)
		// }
	}
	tx.Commit(context.Background())

	return make([]Timeseries, 0), err
}

// func DeleteTimeseries() (something?) {}
