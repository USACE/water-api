package chartserver

import (
	"context"

	"github.com/USACE/water-api/api/helpers"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DamProfileChartInput struct {
	Pool      float64 `querystring:"pool" json:"pool"`
	Tail      float64 `querystring:"tail"`
	Inflow    float64 `querystring:"inflow" json:"inflow"`
	Outflow   float64 `querystring:"outflow" json:"outflow"`
	DamTop    float64 `querystring:"damTop" json:"damtottom" db:"damtop"`
	DamBottom float64 `querystring:"damBottom" json:"dambottom" db:"dambottom"`
	Levels    []Level `querystring:"level" json:"levels" db:"levels"`
}

type Level struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func (s ChartServer) DamProfileChart(input DamProfileChartInput) (string, error) {
	u := *s.URL
	u.Path = u.Path + "/dam-profile-chart"                   // Build URL Path
	u.RawQuery = helpers.StructToQueryValues(input).Encode() // Build URL Query Params
	return helpers.HTTPGet(&u)
}

func GetDamProfileByLocation(db *pgxpool.Pool, locationSlug *string) (*DamProfileChartInput, error) {
	visualizationTypeId, _ := uuid.Parse("53da77d0-6550-4f02-abf8-4bcd1a596a7c")

	var damProfileSQL = `
		WITH lvl_ts AS (
			SELECT 
				vvm.visualization_id,
				vvm.variable AS variable,
				t.datasource_key AS key,
				dt.slug AS datasource_slug,
				t.latest_time,
				t.latest_value 
				FROM visualization_variable_mapping vvm 
				JOIN timeseries t ON t.id = vvm.timeseries_id 
				JOIN datasource d ON d.id = t.datasource_id 
				JOIN datasource_type dt ON dt.id = d.datasource_type_id 
				JOIN visualization v ON v.id = vvm.visualization_id 
			WHERE dt.slug = 'cwms-levels'
		),		
		viz_ts AS (
			SELECT 
				vvm.visualization_id,
				vvm.variable AS variable,
				t.datasource_key AS key,
				dt.slug AS datasource_slug,
				t.latest_time,
				t.latest_value 
				FROM visualization_variable_mapping vvm 
				JOIN timeseries t ON t.id = vvm.timeseries_id 
				JOIN datasource d ON d.id = t.datasource_id 
				JOIN datasource_type dt ON dt.id = d.datasource_type_id 
				JOIN visualization v ON v.id = vvm.visualization_id
			WHERE dt.slug = 'cwms-timeseries'
		)		
		
		SELECT
			(SELECT latest_value FROM viz_ts WHERE variable = 'pool' AND visualization_id = v.id) AS pool,
			(SELECT latest_value FROM viz_ts WHERE variable = 'inflow' AND visualization_id = v.id) AS inflow,
			(SELECT latest_value FROM viz_ts WHERE variable = 'outflow' AND visualization_id = v.id) AS outflow,
			(SELECT latest_value FROM viz_ts WHERE variable = 'tail' AND visualization_id = v.id) AS tail,
			(SELECT latest_value FROM lvl_ts WHERE variable = 'streambed' AND visualization_id = v.id) AS dambottom,
			(SELECT latest_value FROM lvl_ts WHERE variable = 'top-of-dam' AND visualization_id = v.id) AS damtop,
			COALESCE(json_agg(json_build_object(
											'name', lvl_ts.variable,									
											'value', lvl_ts.latest_value
										)), '[]') AS levels
		FROM visualization v
		JOIN visualization_variable_mapping vvm ON vvm.visualization_id = v.id 
		JOIN timeseries t ON t.id = vvm.timeseries_id 
		JOIN datasource d ON d.id = t.datasource_id 
		JOIN datasource_type dt ON dt.id = d.datasource_type_id 
		JOIN "location" l ON l.id = v.location_id 
		LEFT JOIN viz_ts ON viz_ts.visualization_id = vvm.visualization_id AND viz_ts.variable = vvm.variable
		JOIN lvl_ts ON lvl_ts.visualization_id = v.id AND lvl_ts.key = t.datasource_key
		WHERE v.type_id = $2
		AND lower(l.slug) = lower($1)
		GROUP BY 
		v.id
		LIMIT 1`

	var v DamProfileChartInput
	if err := pgxscan.Get(context.Background(), db, &v, damProfileSQL, locationSlug, visualizationTypeId); err != nil {
		return nil, err
	}
	return &v, nil
}
