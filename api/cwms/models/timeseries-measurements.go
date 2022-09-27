package models

import (
	"context"
	"errors"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Measurements struct {
	Times  []time.Time `json:"times,omitempty"`
	Values []float64   `json:"values,omitempty"`
}

func (m Measurements) LatestTime() *time.Time {
	if len(m.Times) > 0 {
		return &m.Times[len(m.Times)-1]
	}
	return nil
}

func (m Measurements) LatestValue() *float64 {
	if len(m.Values) > 0 {
		return &m.Values[len(m.Values)-1]
	}
	return nil
}

func CreateOrUpdateTimeseriesMeasurements(db *pgxpool.Pool, c TimeseriesCollection) ([]Timeseries, error) {

	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]Timeseries, 0), err
	}
	defer tx.Rollback(context.Background())

	updatedIDs := make([]uuid.UUID, 0)

	for _, t := range c.Items {

		// If measurements aren't properly provided, skip this item
		if t.Measurements == nil || t.Measurements.Times == nil || t.Measurements.Values == nil {
			// fmt.Println("skipping " + t.Key)
			continue
		}

		rows, err := tx.Query(
			context.Background(),
			`UPDATE timeseries 
			SET latest_time = $4, latest_value = $5
			WHERE datasource_key = $3
			AND datasource_id = (
				SELECT d.id FROM datasource d 
				JOIN datasource_type dt ON dt.id = d.datasource_type_id 
				JOIN provider p ON p.id = d.provider_id 
				WHERE lower(p.slug) = lower($2) AND lower(dt.slug) = lower($1)
			)
			RETURNING id`,
			t.DatasourceType, t.Provider, t.Key, t.Measurements.LatestTime(), t.Measurements.LatestValue(),
		)
		if err != nil {
			return make([]Timeseries, 0), err
		}
		var id uuid.UUID
		if err := pgxscan.ScanOne(&id, rows); err != nil {
			tx.Rollback(context.Background())
			return c.Items, err
		} else {
			updatedIDs = append(updatedIDs, id)
		}

	}
	tx.Commit(context.Background())

	if len(updatedIDs) == 0 {
		return nil, errors.New("no records updated")
	}

	return nil, nil
}
