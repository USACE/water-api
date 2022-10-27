package models

import (
	"context"
	"time"

	"github.com/USACE/water-api/api/timeseries"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Extract struct for the return of measurements for a watershed
type Extract struct {
	SiteNumber  string      `json:"site_number"`
	Code        string      `json:"code"`
	LocationID  uuid.UUID   `json:"location_id"`
	Name        string      `json:"name"`
	ParameterID uuid.UUID   `json:"paramter_id"`
	Times       []time.Time `json:"times"`
	Values      []float64   `json:"values"`
}

// WatershedExtract
func WatershedExtract(db *pgxpool.Pool, slug string, tw *timeseries.TimeWindow) ([]Extract, error) {
	ext := make([]Extract, 0)
	rows, err := db.Query(context.Background(),
		`SELECT r1.site_number, r1."name", r1.code, r1.location_id, r1.parameter_id, array_agg(r1."time") AS times, array_agg(r1.value) AS values
		FROM
		(
			SELECT us.site_number, us.station_name AS "name", up.code, usp.location_id, usp.parameter_id, um."time", um.value 
			FROM watershed_usgs_sites wus
			JOIN
			(
			SELECT "time", value, usgs_site_parameters_id from usgs_measurements
			WHERE "time" >= $2::timestamptz AND "time" <= $3::timestamptz
			ORDER BY usgs_site_parameters_id, "time" desc
			) AS um on um.usgs_site_parameters_id = wus.usgs_site_parameter_id
		JOIN
		usgs_site_parameters usp 
		ON usgs_site_parameters_id = usp.id
		JOIN usgs_site us 
		ON us.location_id = usp.location_id
		JOIN usgs_parameter up
		ON up.id = usp.parameter_id
		WHERE wus.watershed_id = (SELECT id FROM watershed w WHERE slug = $1)
		ORDER BY us.site_number, up.code, um."time"
		) r1
		GROUP BY r1.site_number, r1."name", r1.code, r1.location_id, r1.parameter_id`,
		slug, tw.After, tw.Before,
	)
	if err != nil {
		return nil, err
	}
	if err = pgxscan.ScanAll(&ext, rows); err != nil {
		return nil, err
	}

	return ext, nil
}
