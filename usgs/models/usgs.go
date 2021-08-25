package models

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/USACE/water-api/helpers"
)

type Parameter struct {
	Id          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
}

type SiteParameter struct {
	UsgsId             string   `json:"usgs_id" db:"usgs_id"`
	UsgsParameterCodes []string `json:"usgs_parameter_codes" db:"usgs_parameter_codes"`
}
type SiteParameterCollection struct {
	Items []SiteParameter `json:"items"`
}

type SiteInfo struct {
	UsgsId            string           `json:"usgs_id" db:"usgs_id"`
	Name              string           `json:"name"`
	Geometry          helpers.Geometry `json:"geometry"`
	Elevation         *float32         `json:"elevation"`
	HorizontalDatumId int              `json:"horizontal_datum_id" db:"horizontal_datum_id"`
	VerticallDatumId  int              `json:"vertical_datum_id" db:"vertical_datum_id"`
	Huc               *string          `json:"huc"`
	StateAbbrev       *string          `json:"state_abbrev" db:"state_abbrev"`
}

type Site struct {
	ID string `json:"-"`
	SiteInfo
	UsgsParameterCodes []string   `json:"parameter_codes" db:"parameter_codes"`
	State              string     `json:"state"`
	VerticalDatum      string     `json:"vertical_datum" db:"vertical_datum"`
	CreateDate         time.Time  `json:"create_date" db:"create_date"`
	UpdateDate         *time.Time `json:"update_date" db:"update_date"`
}

func (s1 Site) IsEquivalent(s2 Site) bool {
	// Compare values that must be set in payload
	return reflect.DeepEqual(s1.SiteInfo, s2.SiteInfo)
}

type SiteFilter struct {
	StateAbbrev *string `json:"state_abbrev" param:"state_abbrev"`
}
type SiteCollection struct {
	Items []Site `json:"items"`
}

func (c *SiteCollection) UnmarshalJSON(b []byte) error {
	switch helpers.JSONType(b) {
	case "ARRAY":
		return json.Unmarshal(b, &c.Items)
	case "OBJECT":
		c.Items = make([]Site, 1)
		return json.Unmarshal(b, &c.Items[0])
	default:
		return errors.New("payload not recognized as JSON array or object")
	}
}

func (c *SiteParameterCollection) UnmarshalJSON(b []byte) error {
	switch helpers.JSONType(b) {
	case "ARRAY":
		return json.Unmarshal(b, &c.Items)
	case "OBJECT":
		c.Items = make([]SiteParameter, 1)
		return json.Unmarshal(b, &c.Items[0])
	default:
		return errors.New("payload not recognized as JSON array or object")
	}
}

func ListSitesQuery(sf *SiteFilter) (sq.SelectBuilder, error) {

	q := sq.Select(`id,
					usgs_id,   
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
	q = q.OrderBy("name")

	// Unfiltered
	return q.PlaceholderFormat(sq.Dollar), nil
}

func ListSites(db *pgxpool.Pool, sf *SiteFilter) ([]Site, error) {

	q, err := ListSitesQuery(sf)
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
			`INSERT INTO usgs_site (usgs_id, name, geometry, elevation, horizontal_datum_id, vertical_datum_id, huc, state_abbrev) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
			m.SiteInfo.UsgsId, m.SiteInfo.Name, m.SiteInfo.Geometry.EWKT(8), m.SiteInfo.Elevation, m.SiteInfo.HorizontalDatumId, m.SiteInfo.VerticallDatumId, m.SiteInfo.Huc, m.SiteInfo.StateAbbrev,
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
	// Where ID In (...)
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
			`UPDATE usgs_site SET name=$2, geometry=$3, elevation=$4, horizontal_datum_id=$5, vertical_datum_id=$6, huc=$7, state_abbrev=$8, update_date=CURRENT_TIMESTAMP WHERE usgs_id = $1 RETURNING id`,
			s.SiteInfo.UsgsId, s.SiteInfo.Name, s.SiteInfo.Geometry.EWKT(8), s.SiteInfo.Elevation, s.SiteInfo.HorizontalDatumId, s.SiteInfo.VerticallDatumId, s.SiteInfo.Huc, s.SiteInfo.StateAbbrev,
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

func ListParameters(db *pgxpool.Pool) ([]Parameter, error) {

	q := sq.Select(`id, code, description`).From("usgs_parameter")
	sql, args, err := q.ToSql()

	if err != nil {
		return make([]Parameter, 0), err
	}
	pp := make([]Parameter, 0)

	if err := pgxscan.Select(context.Background(), db, &pp, sql, args...); err != nil {
		return make([]Parameter, 0), err
	}
	return pp, nil
}

func CreateSiteParameters(db *pgxpool.Pool, ss []SiteParameter) ([]SiteParameter, error) {

	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]SiteParameter, 0), err
	}
	defer tx.Rollback(context.Background())
	for _, m := range ss {
		// fmt.Println(m.UsgsId)
		// fmt.Println(m.UsgsParameterCodes[0])
		rows, err := tx.Query(
			context.Background(),
			`INSERT INTO usgs_site_parameters (usgs_site_id, usgs_parameter_id) 
			VALUES ((select id from usgs_site where usgs_id = $1), (select id from usgs_parameter where code = $2)) RETURNING id`,
			m.UsgsId, m.UsgsParameterCodes[0],
		)
		if err != nil {
			tx.Rollback(context.Background())
			return make([]SiteParameter, 0), err
		}
		var id uuid.UUID
		if err := pgxscan.ScanOne(&id, rows); err != nil {
			tx.Rollback(context.Background())
			return make([]SiteParameter, 0), err
		}
	}
	tx.Commit(context.Background())

	return ss, nil

}

// NOT USED - SAVE FOR WATERSHED/SITE/PARAM Enabled
// func ListParametersEnabled(db *pgxpool.Pool) ([]SiteParameter, error) {

// 	q := sq.Select(`usgs_site_id, usgs_id, usgs_parameter_id, usgs_parameter_code`).From("v_usgs_site_parameters_enabled")

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
