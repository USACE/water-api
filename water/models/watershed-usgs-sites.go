package models

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type WatershedSiteParameter struct {
	WatershedSlug string `param:"watershed_slug"`
	SiteNumber    string `param:"site_number"`
	ParameterCode string `param:"parameter_code"`
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
		WHERE watershed_id = (select id from watershed where slug = $1 limit 1)
		AND usgs_site_id = (select id from usgs_site where site_number = $2 limit 1)
		AND usgs_parameter_id = (select id from usgs_parameter where code = $3 limit 1)
		`, w.WatershedSlug, w.SiteNumber, w.ParameterCode,
	); err != nil {
		return err
	}
	return nil

}
