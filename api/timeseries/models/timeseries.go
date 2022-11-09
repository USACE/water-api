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
	TimeseriesFilter struct {
		IDs              *[]uuid.UUID // intentionally not supported as query param
		Datatype         *string      `query:"datatype"`
		Provider         *string      `query:"provider"`
		EtlValuesEnabled *bool        `query:"etl_values_enabled"`
		Q                *string      `query:"q"`
	}

	Timeseries struct {
		Provider         string           `json:"provider"`
		ProviderName     string           `json:"provider_name"`
		Datatype         string           `json:"datatype"`
		DatatypeName     string           `json:"datatype_name"`
		Key              string           `json:"key"`
		Location         *string          `json:"location"`                       // location slug
		LatestValue      *[]interface{}   `json:"latest_value" db:"latest_value"` // e.g. ["2022-09-27T12:00:00-05:00", 888.14]
		Values           *[][]interface{} `json:"values,omitempty"`               // may be empty [] or [["2022-09-27T12:00:00-05:00", 888.00], ["2022-09-27T13:00:00-05:00", 888.15]]
		EtlValuesEnabled *bool            `json:"etl_values_enabled,omitempty" db:"etl_values_enabled"`
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
		`p.slug  			    AS provider,
		 p.name                 AS provider_name,
		 dt.slug 			    AS datatype,
		 dt.name                AS datatype_name,
		 l.slug				    AS location,
		 t.datasource_key 	    AS key,
		 json_build_array(
			t.latest_time,
			t.latest_value
		)::json                 AS latest_value`,
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

		// Filter by etl_values_enabled
		if f.EtlValuesEnabled != nil {
			q = q.Where("t.etl_values_enabled = ?", f.EtlValuesEnabled)
		}

		// Filter by list of known UUIDs
		// This is used after CreateTimeseries when UUIDs of newly created locations are known
		if f.IDs != nil {
			q = q.Where(sq.Eq{"t.id": f.IDs})
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

// func CreateTimeseries() ([]Timeseries, error) {}
func (tsc TimeseriesCollection) Create(db *pgxpool.Pool, providerSlug string) ([]Timeseries, error) {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]Timeseries, 0), err
	}
	defer tx.Rollback(context.Background())

	newIDs := make([]uuid.UUID, 0)
	for _, t := range tsc.Items {
		rows, err := tx.Query(
			context.Background(),
			`INSERT INTO timeseries (datasource_id, datasource_key, location_id)
             VALUES (
                (SELECT id FROM v_datasource WHERE datatype = LOWER($1) AND provider = LOWER($2)),
                $3,
                (SELECT id FROM v_location WHERE slug = LOWER($4) AND provider = LOWER($2))
             ) ON CONFLICT ON CONSTRAINT timeseries_unique_datasource DO NOTHING
			 RETURNING id`, t.Datatype, t.Provider, t.Key, t.Location,
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
