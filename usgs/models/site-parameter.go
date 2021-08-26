package models

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

// SiteParameter is a relationship of a ParameterCode to a SiteNumber
type SiteParameter struct {
	SiteNumber     string   `json:"site_number" db:"site_number"`
	ParameterCodes []string `json:"parameter_codes" db:"parameter_codes"`
}

func CreateSiteParameters(db *pgxpool.Pool, ss []SiteParameter) ([]SiteParameter, error) {

	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]SiteParameter, 0), err
	}
	defer tx.Rollback(context.Background())
	for _, m := range ss {
		rows, err := tx.Query(
			context.Background(),
			`INSERT INTO usgs_site_parameters (site_uid, parameter_uid) VALUES
			((select uid from usgs_site where site_number = $1), (select uid from usgs_parameter where code = $2))
			RETURNING uid`, m.SiteNumber, m.ParameterCodes[0],
		)
		if err != nil {
			tx.Rollback(context.Background())
			return make([]SiteParameter, 0), err
		}
		var uid uuid.UUID
		if err := pgxscan.ScanOne(&uid, rows); err != nil {
			tx.Rollback(context.Background())
			return make([]SiteParameter, 0), err
		}
	}
	tx.Commit(context.Background())

	return ss, nil

}
