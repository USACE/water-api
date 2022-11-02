package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/exp/slices"
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

	timeseriesIdsUpdated := make([]uuid.UUID, 0)

	for _, t := range c.Items {

		// If measurements aren't properly provided, skip this item
		if t.Measurements == nil || t.Measurements.Times == nil || t.Measurements.Values == nil {
			// fmt.Println("skipping " + t.Key)
			continue
		}
		/********************************************************
		Get the timeseries_id based on the current item's payload
		(datatype, provider, key)
		********************************************************/
		var timeseriesId uuid.UUID
		if err := pgxscan.Get(
			context.Background(), db, &timeseriesId, ` 
			SELECT t.id FROM timeseries t 
			JOIN datasource d ON d.id = t.datasource_id 
			JOIN datatype dt ON dt.id = d.datatype_id
			JOIN provider p ON p.id = d.provider_id 
			WHERE lower(dt.slug) = lower($1)
			AND lower(p.slug) = lower($2)
			AND t.datasource_key = $3
			`, t.Datatype, t.Provider, t.Key,
		); err != nil {
			return nil, err
		}

		/********************************************************
		Iterate over the measurements and insert into db
		********************************************************/

		for idx := range t.Measurements.Times {

			rows, err := tx.Query(
				context.Background(),
				`INSERT INTO timeseries_measurement (timeseries_id, time, value) VALUES ($1, $2, $3)
				ON CONFLICT ON CONSTRAINT timeseries_id_unique_time
				DO UPDATE SET value = EXCLUDED.value
				RETURNING timeseries_id`,
				timeseriesId, t.Measurements.Times[idx], t.Measurements.Values[idx],
			)

			if err != nil {
				tx.Rollback(context.Background())
				return make([]Timeseries, 0), err
			}

			var id uuid.UUID
			if err := pgxscan.ScanOne(&id, rows); err != nil {
				tx.Rollback(context.Background())
				return make([]Timeseries, 0), err

			} else {
				if !slices.Contains(timeseriesIdsUpdated, id) {
					timeseriesIdsUpdated = append(timeseriesIdsUpdated, id)
				}

			}
		}
	}
	tx.Commit(context.Background())

	/********************************************************
	After the above transaction has been completed...
	Iterate over the timeseriesIdsUpdated and update the
	lastest_time and latest_value in the timeseries table
	********************************************************/

	tx, err = db.Begin(context.Background())
	if err != nil {
		return make([]Timeseries, 0), err
	}
	defer tx.Rollback(context.Background())

	for i := range timeseriesIdsUpdated {

		timeseriesId := timeseriesIdsUpdated[i]
		fmt.Println(timeseriesId)

		type Row struct {
			T time.Time
		}
		var maxTime Row
		if err := pgxscan.Get(
			context.Background(), db, &maxTime, `
			SELECT max(time) as t FROM timeseries_measurement WHERE timeseries_id = $1
			`, timeseriesId,
		); err != nil {
			println(err.Error())
			return nil, err
		}

		fmt.Println(maxTime)

		_, err := tx.Exec(
			context.Background(),
			`UPDATE timeseries 
			SET 
				latest_time = $2, 
				latest_value = (SELECT value from timeseries_measurement where timeseries_id = $1 and time = $2)
			WHERE id = $1`,
			timeseriesId, maxTime.T,
		)

		if err != nil {
			tx.Rollback(context.Background())
			return make([]Timeseries, 0), err
		}
	}

	tx.Commit(context.Background())

	if len(timeseriesIdsUpdated) == 0 {
		return nil, errors.New("no records updated")
	}

	return nil, nil
}
