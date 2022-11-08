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

type (
	TimeseriesFilter struct {
		IDs      *[]uuid.UUID // intentionally not supported as query param
		Datatype *string      `query:"datatype"`
		Provider *string      `query:"provider"`
		Q        *string      `query:"q"`
	}

	Timeseries struct {
		Provider     string        `json:"provider"`
		ProviderName string        `json:"provider_name"`
		Datatype     string        `json:"datatype"`
		DatatypeName string        `json:"datatype_name"`
		Key          string        `json:"key"`
		Location     *string       `json:"location"`                         // location slug
		LocationCode *string       `json:"location_code" db:"location_code"` // location code
		LatestTime   *time.Time    `json:"latest_time" db:"latest_time"`
		LatestValue  *float64      `json:"latest_value" db:"latest_value"`
		Measurements *Measurements `json:"measurements,omitempty"` // may be empty
		//	Creates a Timeseries may or may not have a Measurements struct.
		// A measurements struct has two fields. Times, Values. Each is an array with zero or more values []
		//	Creates metadata that will be associated with all related timeseries measurements
	}

	TimeseriesCollection struct {
		Items []Timeseries `json:"items"`
	}
)

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

	q := sq.Select(
		`p.slug  			AS provider,
		 p.name             AS provider_name,
		 dt.slug 			AS datatype,
		 dt.name            AS datatype_name,
		 l.code             AS location_code,
		 l.slug				AS location,
		 t.datasource_key 	AS key,
		 t.latest_time,
		 t.latest_value`,
	).From(
		"timeseries t",
	)

	jDS, jDSParams := "datasource ds ON ds.id  = t.datasource_id", make([]interface{}, 0) // join datasource
	jP, jPParams := "provider     p  ON p.id   = ds.provider_id", make([]interface{}, 0)  // join provider
	jDT, jDTParams := "datatype   dt ON dt.id  = ds.datatype_id", make([]interface{}, 0)  // join datatype

	if f != nil {

		// Filter by Provider
		if f.Provider != nil {
			// Limit datasource join by ?provider= query param
			jDS += " AND ds.provider_id = (SELECT id from provider WHERE slug = LOWER(?))"
			jDSParams = append(jDSParams, f.Provider)
			// Limit provider join by ?provider= query param
			jP += " AND p.slug = LOWER(?)"
			jPParams = append(jPParams, f.Provider)
			// WHERE
			q = q.Where("p.slug = LOWER(?)", f.Provider)
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
			q = q.Where("dt.slug = LOWER(?)", f.Datatype)
		}

		// Filter by search string
		if f.Q != nil {
			q = q.Where("t.datasource_key ILIKE '%' || lower(?) || '%' ", f.Q)
		}
	}

	q = q.Join(jDS, jDSParams...)                        // join datasource
	q = q.Join(jP, jPParams...)                          // join provider
	q = q.Join(jDT, jDTParams...)                        // join datatype
	q = q.LeftJoin("location l on l.id = t.location_id") // left join location

	q = q.OrderBy("p.slug, t.datasource_key")

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

// func GetTimeseries() (*Timeseries, error) {}

// func CreateTimeseries() ([]Timeseries, error) {}
func CreateTimeseries(db *pgxpool.Pool, c TimeseriesCollection) ([]Timeseries, error) {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]Timeseries, 0), err
	}
	defer tx.Rollback(context.Background())

	newIDs := make([]uuid.UUID, 0)
	for _, t := range c.Items {
		rows, err := tx.Query(
			context.Background(),
			`INSERT INTO timeseries (datasource_id, datasource_key, location_id)
			VALUES(
				(
					SELECT id
					  FROM datasource
					 WHERE datatype_id = (SELECT id FROM datatype WHERE slug = LOWER($1))
					   AND provider_id = (SELECT id FROM provider WHERE slug = LOWER($2))
				),
				$3,
				(
					SELECT id
					  FROM location
					 WHERE provider_id = (SELECT id FROM provider WHERE slug = LOWER($2))
					   AND code        = LOWER($4)
				)
			ON CONFLICT DO NOTHING		
			RETURNING id`,
			t.Datatype, t.Provider, t.Key, t.LocationCode,
		)
		if err != nil {
			tx.Rollback(context.Background())
			return make([]Timeseries, 0), err
		}
		var id uuid.UUID
		if err := pgxscan.ScanOne(&id, rows); err != nil {
			continue
		}
		newIDs = append(newIDs, id)
	}
	tx.Commit(context.Background())

	return ListTimeseries(db, &TimeseriesFilter{IDs: &newIDs})
}

// func DeleteTimeseries() (something?) {}
