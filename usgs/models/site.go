package models

import (
	"context"
	"reflect"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/USACE/water-api/helpers"
)

type Site struct {
	UID            uuid.UUID  `json:"-" db:"id"`
	ParameterCodes []string   `json:"parameter_codes" db:"parameter_codes"`
	State          string     `json:"state"`
	VerticalDatum  string     `json:"vertical_datum" db:"vertical_datum"`
	CreateDate     time.Time  `json:"create_date" db:"create_date"`
	UpdateDate     *time.Time `json:"update_date" db:"update_date"`
	SiteInfo
}

func (s1 Site) IsEquivalent(s2 Site) bool {
	// Compare values that must be set in payload
	return reflect.DeepEqual(s1.SiteInfo, s2.SiteInfo)
}

type SiteInfo struct {
	SiteNumber        string           `json:"site_number" db:"site_number"`
	Name              string           `json:"name"`
	Geometry          helpers.Geometry `json:"geometry"`
	Elevation         *float32         `json:"elevation"`
	HorizontalDatumId int              `json:"horizontal_datum_id" db:"horizontal_datum_id"`
	VerticallDatumId  int              `json:"vertical_datum_id" db:"vertical_datum_id"`
	Huc               *string          `json:"huc"`
	StateAbbrev       *string          `json:"state_abbrev" db:"state_abbrev"`
}

type SiteFilter struct {
	StateAbbrev *string `json:"state" query:"state"`
	Q           *string `query:"q"`
}

func ListSitesQuery(sf *SiteFilter) (sq.SelectBuilder, error) {

	q := sq.Select(`id,
					site_number,   
					name,		            
		            geometry,
		            elevation,
					horizontal_datum_id,
					vertical_datum_id,
					vertical_datum,
					huc,
					state_abbrev,
					parameter_codes,
					state,
					create_date,
					update_date`).From("v_usgs_site")

	if sf != nil {
		// Filter by StateID
		if sf.StateAbbrev != nil {
			q = q.Where("state_abbrev = ?", strings.ToUpper(*sf.StateAbbrev))
		}
	}

	// Unfiltered
	return q.PlaceholderFormat(sq.Dollar), nil
}

func SearchSites(db *pgxpool.Pool, f *SiteFilter) ([]Site, error) {
	q, err := ListSitesQuery(f)
	if err != nil {
		return make([]Site, 0), err
	}
	// Filter by Query String
	q = q.Where("name ILIKE '%' || ? || '%' ORDER BY name LIMIT 10", f.Q)
	sql, args, err := q.ToSql()
	if err != nil {
		return make([]Site, 0), err
	}
	ss := make([]Site, 0)
	if err := pgxscan.Select(context.Background(), db, &ss, sql, args...); err != nil {
		return make([]Site, 0), err
	}
	return ss, nil
}

func ListSites(db *pgxpool.Pool, sf *SiteFilter) ([]Site, error) {

	q, err := ListSitesQuery(sf)
	q = q.OrderBy("name")
	if err != nil {
		return make([]Site, 0), err
	}
	sql, args, err := q.ToSql()
	if err != nil {
		return make([]Site, 0), err
	}
	ss := make([]Site, 0)
	if err := pgxscan.Select(context.Background(), db, &ss, sql, args...); err != nil {
		return make([]Site, 0), err
	}
	return ss, nil
}

func ListSitesForIDs(db *pgxpool.Pool, IDs []uuid.UUID) ([]Site, error) {
	// Base Locations Query
	q, err := ListSitesQuery(nil)
	if err != nil {
		return make([]Site, 0), err
	}
	// Where ID In (...)
	q = q.Where(sq.Eq{"id": IDs})
	sql, args, err := q.ToSql()
	if err != nil {
		return make([]Site, 0), err
	}
	ss := make([]Site, 0)
	if err := pgxscan.Select(context.Background(), db, &ss, sql, args...); err != nil {
		return make([]Site, 0), err
	}
	return ss, nil
}

func CreateSites(db *pgxpool.Pool, nn []Site) ([]Site, error) {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]Site, 0), err
	}
	defer tx.Rollback(context.Background())
	newIDs := make([]uuid.UUID, 0)
	for _, m := range nn {
		rows, err := tx.Query(
			context.Background(),
			`INSERT INTO usgs_site (site_number, name, geometry, elevation, horizontal_datum_id, vertical_datum_id, huc, state_abbrev) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id`,
			m.SiteInfo.SiteNumber, m.SiteInfo.Name, m.SiteInfo.Geometry.EWKT(8), m.SiteInfo.Elevation, m.SiteInfo.HorizontalDatumId, m.SiteInfo.VerticallDatumId, m.SiteInfo.Huc, m.SiteInfo.StateAbbrev,
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
	tx.Commit(context.Background())

	return ListSitesForIDs(db, newIDs)
}

func GetSiteByID(db *pgxpool.Pool, siteID *uuid.UUID) (*Site, error) {
	// Base Sites Query
	q, err := ListSitesQuery(nil)
	if err != nil {
		return nil, err
	}
	// Where UID In (...)
	q = q.Where("id = ?", siteID)
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	var s Site
	if err := pgxscan.Get(context.Background(), db, &s, sql, args...); err != nil {
		return nil, err
	}
	return &s, nil
}

func UpdateSites(db *pgxpool.Pool, nn []Site) ([]Site, error) {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]Site, 0), err
	}
	defer tx.Rollback(context.Background())
	newIDs := make([]uuid.UUID, 0)
	for _, s := range nn {
		rows, err := tx.Query(
			context.Background(),
			`UPDATE usgs_site SET name=$2, geometry=$3, elevation=$4, horizontal_datum_id=$5, vertical_datum_id=$6, huc=$7, state_abbrev=$8, update_date=CURRENT_TIMESTAMP WHERE site_number = $1
			 RETURNING id`,
			s.SiteInfo.SiteNumber, s.SiteInfo.Name, s.SiteInfo.Geometry.EWKT(8), s.SiteInfo.Elevation, s.SiteInfo.HorizontalDatumId, s.SiteInfo.VerticallDatumId, s.SiteInfo.Huc, s.SiteInfo.StateAbbrev,
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
	tx.Commit(context.Background())

	return ListSitesForIDs(db, newIDs)
}

// NOT USED - SAVE FOR WATERSHED/SITE/PARAM Enabled
// func ListParametersEnabled(db *pgxpool.Pool) ([]SiteParameter, error) {

// 	q := sq.Select(`usgs_site_id, site_number, usgs_parameter_id, usgs_parameter_code`).From("v_usgs_site_parameters_enabled")

// 	sql, args, err := q.ToSql()

// 	if err != nil {
// 		return make([]SiteParameter, 0), err
// 	}
// 	pp := make([]SiteParameter, 0)

// 	if err := pgxscan.Select(context.Background(), db, &pp, sql, args...); err != nil {
// 		return make([]SiteParameter, 0), err
// 	}
// 	return pp, nil
// }
