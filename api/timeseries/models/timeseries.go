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
		LatestValue      *[]interface{}   `json:"latest_value" db:"latest_value"` // e.g. ["2022-09-27T12:00:00-05:00", 888.14]
		Values           *[][]interface{} `json:"values,omitempty"`               // may be empty [] or [["2022-09-27T12:00:00-05:00", 888.00], ["2022-09-27T13:00:00-05:00", 888.15]]
		EtlValuesEnabled *bool            `json:"etl_values_enabled,omitempty" db:"etl_values_enabled"`
		Location         struct {
			Slug     *string `json:"slug"`     // Optional Location Information; required to establish linkage to unique location on Create
			Provider *string `json:"provider"` // Optional Location Information; required to establish linkage to unique location on Create
			Datatype *string `json:"datatype"` // Optional Location Information; required to establish linkage to unique location on Create
			Code     *string `json:"code"`     // Optional Location Information; required to establish linkage to unique location on Create
		} `json:"location"` // todo; consider using a fully populated `location.LocationInfo` struct here
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
		`provider, provider_name, datatype, datatype_name, key, latest_value, location`,
	).From(
		"v_timeseries t",
	)

	if f != nil {
		// Filter by Provider
		if f.Provider != nil {
			q = q.Where("provider = LOWER(?)", f.Provider)
		}

		// Filter by Datatype
		if f.Datatype != nil {
			q = q.Where("datatype = LOWER(?)", f.Datatype)
		}

		// Filter by search string
		if f.Q != nil {
			q = q.Where("key ILIKE '%' || lower(?) || '%' ", f.Q)
		}

		// Filter by etl_values_enabled
		if f.EtlValuesEnabled != nil {
			q = q.Where("etl_values_enabled = ?", f.EtlValuesEnabled)
		}

		// Filter by list of known UUIDs
		// This is used after CreateTimeseries when UUIDs of newly created locations are known
		if f.IDs != nil {
			q = q.Where(sq.Eq{"id": f.IDs})
		}
	}

	q = q.OrderBy("provider, key")

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
                (SELECT id FROM v_location WHERE code = LOWER($4) AND provider = LOWER($$5))
             ) ON CONFLICT ON CONSTRAINT timeseries_unique_datasource DO NOTHING
			 RETURNING id`, t.Datatype, t.Provider, t.Key, t.Location.Code, t.Location.Provider,
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
