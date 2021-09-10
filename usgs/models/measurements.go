package models

import (
	"context"
	"time"

	"github.com/USACE/water-api/timeseries"
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
func CreateOrUpdateMeasurements(db *pgxpool.Pool, c ParameterMeasurementCollection) (map[string]string, error) {
	// Loop through the array of parameter measurements
	tx, err := db.Begin(context.Background())
	if err != nil {
		return make(map[string]string), err
	}
	defer tx.Rollback(context.Background())
	// newIDs := make([]uuid.UUID, 0)
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
				return make(map[string]string), err
			}
			var id uuid.UUID
			if err := pgxscan.ScanOne(&id, rows); err != nil {
				tx.Rollback(context.Background())
				return make(map[string]string), err
				// } else {
				// 	newIDs = append(newIDs, id)
			}
		}
	}
	tx.Commit(context.Background())

	// return ListSitesForIDs(db, newIDs)
	return make(map[string]string), nil
}

// ListMeasurements returns time and value for the USGS location
// filtered by a time range.
func ListUSGSMeasurements(db *pgxpool.Pool, site_number *string, parameters []string, tw *timeseries.TimeWindow) (ParameterMeasurementCollection, error) {
	tx, err := db.Begin(context.Background())
	var pc ParameterMeasurementCollection
	var pm ParameterMeasurements
	if err != nil {
		return pc, err
	}
	defer tx.Rollback(context.Background())
	if len(parameters) == 0 {
		rows, _ := tx.Query(
			context.Background(),
			`WITH s_id AS (
				SELECT s.* FROM a2w_cwms.usgs_site AS s WHERE s.site_number = $1
			), s_parameters AS (
				SELECT sp.* FROM a2w_cwms.usgs_site_parameters AS sp
				WHERE sp.site_id = (SELECT id FROM s_id)
			)
			SELECT p.code FROM a2w_cwms.usgs_parameter p, s_parameters s
			WHERE p.id = s.parameter_id`,
			site_number,
		)
		pgxscan.ScanAll(&parameters, rows)
	}
	for _, parameter := range parameters {
		rows, err := tx.Query(
			context.Background(),
			`WITH s_id AS (
			SELECT id FROM usgs_site s WHERE s.site_number = $1
			), p_id AS (
				SELECT id FROM usgs_parameter p WHERE p.code = $2
			), site_parameter_id AS (
				SELECT id FROM usgs_site_parameters sp
				WHERE parameter_id = (SELECT * FROM p_id) AND site_id = (SELECT * FROM s_id)
			)
			SELECT m.time, m.value
			FROM usgs_measurements m
			WHERE usgs_site_parameters_id = (SELECT * FROM site_parameter_id)
			AND time >= $3 AND time <= $4
			ORDER BY time ASC`,
			site_number,
			parameter,
			tw.After.Format(time.RFC3339),
			tw.Before.Format(time.RFC3339),
		)
		if err != nil {
			tx.Rollback(context.Background())
			return pc, err
		}
		// ms := make([]Measurement, 0)
		var ms []Measurement
		if err := pgxscan.ScanAll(&ms, rows); err != nil {
			tx.Rollback(context.Background())
			return pc, err
		}
		pm.ParameterCode = parameter
		pm.Measurements.Items = ms
		pc.SiteNumber = *site_number
		pc.Items = append(pc.Items, pm)
	}
	return pc, nil
}
