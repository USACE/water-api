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
		json_agg(json_build_object('sites', sites, 'state_abbrev', sites.state_abbrev))
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
			JOIN usgs_parameter p ON
				p.id = w.usgs_parameter_id
			JOIN usgs_site s ON
				s.id = w.usgs_site_id
			LEFT JOIN (
				SELECT
					array_agg(code) AS parameter_codes,
					b.usgs_site_id
				FROM
					usgs_parameter a
				JOIN watershed_usgs_sites b ON
					b.usgs_parameter_id = a.id
				GROUP BY
					b.usgs_site_id
				) code_agg ON
				code_agg.usgs_site_id = w.usgs_site_id
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
		`INSERT INTO watershed_usgs_sites (watershed_id, usgs_site_id, usgs_parameter_id) VALUES
		((select id from watershed where slug = $1), 
		(select id from usgs_site where site_number = $2),
		(select id from usgs_parameter where code = $3))
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
		AND usgs_site_id = (select id from usgs_site where site_number = $2)
		AND usgs_parameter_id = (select id from usgs_parameter where code = $3)
		`, w.WatershedSlug, w.SiteNumber, w.ParameterCode,
	); err != nil {
		return err
	}
	return nil

}
