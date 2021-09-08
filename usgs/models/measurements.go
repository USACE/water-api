package models

import (
	"context"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Measurement is a single object with time and value.
// An array, or collection, of Measurements is a MeasurementCollection
// defined in the MeasurementCollection in the measurements-collection.go file.
type Measurement struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`
}

// ParameterMeasurement defines a single parameter with measurement(s)
// where the measurements are in an array, Measurement struct.
// An array, or collection, of parameter measurements is a ParameterMeasurementCollection
// defined in the measurements-collection.go file.
type ParameterMeasurements struct {
	ParameterCode string                `json:"code" db:"paramter_code"`
	Measurements  MeasurementCollection `json:"measurements"`
}

// CreateOrUpdateTimeseriesMeasurements
func CreateOrUpdateMeasurements(db *pgxpool.Pool, c ParameterMeasurementCollection) ([]Site, error) {
	// Loop through the array of parameter measurements
	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]Site, 0), err
	}
	defer tx.Rollback(context.Background())
	newIDs := make([]uuid.UUID, 0)
	s_number := c.SiteNumber
	pm := c.Items
	for idx := range pm {
		p_code := pm[idx].ParameterCode
		ms := pm[idx].Measurements.Items
		for ndx := range ms {
			m := ms[ndx]
			rows, err := tx.Query(
				context.Background(),
				`WITH s_id AS (
					SELECT id FROM a2w_cwms.usgs_site s WHERE s.site_number = $1
				), p_id AS (
					SELECT id FROM a2w_cwms.usgs_parameter p WHERE p.code = $2
				), site_parameter_id AS (
					SELECT id FROM a2w_cwms.usgs_site_parameters sp
					WHERE parameter_id = (SELECT * FROM p_id) AND site_id = (SELECT * FROM s_id)
				)
				INSERT INTO a2w_cwms.usgs_measurements (time, value, usgs_site_parameters_id) VALUES
				($3, $4, (SELECT * FROM site_parameter_id))
				ON CONFLICT ON CONSTRAINT site_parameters_unique_time
				DO UPDATE SET value = EXCLUDED.value
				RETURNING (SELECT * FROM s_id)`,
				s_number,
				p_code,
				m.Time,
				m.Value,
			)
			if err != nil {
				tx.Rollback(context.Background())
				return make([]Site, 0), err
			}
			var id uuid.UUID
			if err := pgxscan.ScanOne(&id, rows); err != nil {
				tx.Rollback(context.Background())
				return make([]Site, 0), err
			} else {
				newIDs = append(newIDs, id)
			}
		}
	}
	tx.Commit(context.Background())

	return ListSitesForIDs(db, newIDs)
}
