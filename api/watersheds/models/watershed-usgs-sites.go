package models

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

type WatershedSiteParameter struct {
	WatershedSlug string `param:"watershed_slug"`
	SiteNumber    string `param:"site_number"`
	ParameterCode string `param:"parameter_code"`
}

func ListWatershedSiteParameters(db *pgxpool.Pool) ([]byte, error) {

	var b []byte
	err := pgxscan.Get(context.Background(), db, &b,
		`SELECT
		COALESCE(json_agg(json_build_object('sites', sites, 'state_abbrev', sites.state_abbrev)), '[]')
	FROM
		(
		SELECT
			json_agg(t) AS sites,
			t.state_abbrev
		FROM
			(
			SELECT
				DISTINCT
				s.site_number AS site_number,
				s.state_abbrev AS state_abbrev,
				COALESCE(code_agg.parameter_codes, '{}') AS parameter_codes
			FROM
				watershed_usgs_sites w				
			JOIN usgs_site_parameters usp ON
				usp.id = w.usgs_site_parameter_id			
			JOIN usgs_parameter p ON
				p.id = usp.parameter_id
			JOIN usgs_site s ON
				s.id = usp.site_id
			LEFT JOIN (			
				SELECT
					array_agg(DISTINCT code) AS parameter_codes,
					usp.site_id
				FROM
					usgs_parameter a
				JOIN usgs_site_parameters usp ON
					usp.parameter_id = a.id
				JOIN watershed_usgs_sites b ON
					b.usgs_site_parameter_id = usp.id
				GROUP BY
					usp.site_id					
					) code_agg ON
				code_agg.site_id = usp.site_id
		) AS t
		GROUP BY
			t.state_abbrev
		)AS sites
	`)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func CreateWatershedSiteParameter(db *pgxpool.Pool, w *WatershedSiteParameter) error {

	// var wsp WatershedSiteParameter
	if _, err := db.Exec(
		context.Background(),
		`INSERT INTO watershed_usgs_sites (watershed_id, usgs_site_parameter_id) VALUES
		((select id from watershed where slug = $1), 
		(
			SELECT usp.id FROM usgs_site_parameters usp
			JOIN usgs_site us ON us.id = usp.site_id 
			JOIN usgs_parameter up ON up.id = usp.parameter_id 
			WHERE us.site_number = $2
			AND up.code = $3)
		)
		`, w.WatershedSlug, w.SiteNumber, w.ParameterCode,
	); err != nil {
		return err
	}

	return nil
}

func DeleteWatershedSiteParameter(db *pgxpool.Pool, w *WatershedSiteParameter) error {

	if _, err := db.Exec(
		context.Background(),
		`DELETE FROM watershed_usgs_sites
		WHERE watershed_id = (select id from watershed where slug = $1)
		AND usgs_site_parameter_id = (
			SELECT usp.id FROM usgs_site_parameters usp
			JOIN usgs_site us ON us.id = usp.site_id 
			JOIN usgs_parameter up ON up.id = usp.parameter_id 
			WHERE us.site_number = $2
			AND up.code = $3)
		`, w.WatershedSlug, w.SiteNumber, w.ParameterCode,
	); err != nil {
		return err
	}
	return nil

}
