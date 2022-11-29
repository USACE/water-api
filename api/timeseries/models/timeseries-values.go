package models

import (
	"context"
	"errors"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

func (tsc TimeseriesCollection) CreateOrUpdateTimeseriesValues(db *pgxpool.Pool) ([]Timeseries, error) {

	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]Timeseries, 0), err
	}
	defer tx.Rollback(context.Background())

	timeseriesIdsUpdated := make(map[uuid.UUID]bool)

	for _, t := range tsc.Items {

		// If values aren't properly provided, skip this item
		if t.Values == nil {
			continue
		}

		/********************************************************
		Get the timeseries_id based on the current item's payload
		(datatype, provider, key)
		********************************************************/
		var timeseriesId uuid.UUID
		if err := pgxscan.Get(
			context.Background(), db, &timeseriesId, `
			SELECT id
			  FROM v_timeseries
			 WHERE datatype   = LOWER($1)
			   AND provider   = LOWER($2)
			   AND LOWER(key) = LOWER($3)
			`, t.Datatype, t.Provider, t.Key,
		); err != nil {
			continue // move on to the next timeseries
		}

		/********************************************************
		Iterate over the values and insert into db
		********************************************************/
		for _, val := range *t.Values {

			rows, err := tx.Query(
				context.Background(),
				`INSERT INTO timeseries_value (timeseries_id, time, value)
				 VALUES ($1, $2, $3)
				 ON CONFLICT ON CONSTRAINT timeseries_id_unique_time DO UPDATE SET value = EXCLUDED.value
				 RETURNING timeseries_id`, timeseriesId, val[0], val[1],
			)
			if err != nil {
				tx.Rollback(context.Background())
				return make([]Timeseries, 0), err
			}

			var id uuid.UUID
			if err := pgxscan.ScanOne(&id, rows); err != nil {
				tx.Rollback(context.Background())
				return make([]Timeseries, 0), err

			}
			timeseriesIdsUpdated[id] = true
		}
	}
	tx.Commit(context.Background())

	/**********************************************************
	After the above transaction has been committed, so new
	values are available to SELECT and consider in MAX(time)...
	Iterate over the timeseriesIdsUpdated and update the
	lastest_time and latest_value in the timeseries table
	**********************************************************/
	tx, err = db.Begin(context.Background())
	if err != nil {
		return make([]Timeseries, 0), err
	}
	defer tx.Rollback(context.Background())

	for tid := range timeseriesIdsUpdated {
		var row struct{ Time time.Time }
		if err := pgxscan.Get(
			context.Background(), db, &row, `
			SELECT MAX(time) AS time FROM timeseries_value WHERE timeseries_id = $1`, tid,
		); err != nil {
			return make([]Timeseries, 0), err
		}

		_, err := tx.Exec(
			context.Background(),
			`UPDATE timeseries
			 SET latest_time  = $2,
				 latest_value = (
					SELECT value
					  FROM timeseries_value
					 WHERE timeseries_id = $1
				       AND time = $2
				 )
			 WHERE id = $1 AND (latest_time IS NULL OR latest_time <= $2)`,
			tid, row.Time,
		)

		if err != nil {
			tx.Rollback(context.Background())
			return make([]Timeseries, 0), err
		}
	}

	tx.Commit(context.Background())

	if len(timeseriesIdsUpdated) == 0 {
		return make([]Timeseries, 0), errors.New("no records updated")
	}

	// Convert uuids of timeseries that have been updated to a slice
	ids := make([]uuid.UUID, len(timeseriesIdsUpdated))
	for id := range timeseriesIdsUpdated {
		ids = append(ids, id)
	}

	return ListTimeseries(db, &TimeseriesFilter{IDs: &ids})
}
