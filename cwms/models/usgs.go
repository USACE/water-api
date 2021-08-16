package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SiteInfo struct {
	UsgsId            string   `json:"usgs_id" db:"usgs_id"`
	Name              string   `json:"name"`
	Geometry          Geometry `json:"geometry"`
	Elevation         *float32 `json:"elevation"`
	HorizontalDatumId int      `json:"horizontal_datum_id" db:"horizontal_datum_id"`
	VerticallDatumId  int      `json:"vertical_datum_id" db:"vertical_datum_id"`
	Huc               *string  `json:"huc"`
	StateAbbrev       *string  `json:"state_abbrev" db:"state_abbrev"`
}

type Site struct {
	ID string `json:"id"`
	SiteInfo
	State         string     `json:"state"`
	VerticalDatum string     `json:"vertical_datum" db:"vertical_datum"`
	CreateDate    time.Time  `json:"create_date" db:"create_date"`
	UpdateDate    *time.Time `json:"update_date" db:"update_date"`
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
	switch JSONType(b) {
	case "ARRAY":
		return json.Unmarshal(b, &c.Items)
	case "OBJECT":
		c.Items = make([]Site, 1)
		return json.Unmarshal(b, &c.Items[0])
	default:
		return errors.New("payload not recognized as JSON array or object")
	}
}

func ListSitesQuery(sf *SiteFilter) (sq.SelectBuilder, error) {

	q := sq.Select(`s.id,
					s.usgs_id,   
					s.name,		            
		            ST_AsGeoJSON(s.geometry)::json AS geometry,
		            s.elevation,
					s.horizontal_datum_id,
					s.vertical_datum_id,
					v.name AS vertical_datum,
					s.huc,
					s.state_abbrev,
					st.name as state,
					s.create_date,
					s.update_date`,
	).From("usgs_site s").Join("vertical_datum v on v.id=s.vertical_datum_id").Join("tiger_data.state_all st on st.stusps=s.state_abbrev")

	if sf != nil {
		// Filter by StateID
		if sf.StateAbbrev != nil {
			q = q.Where("s.state_abbrev = ?", strings.ToUpper(*sf.StateAbbrev))
		}
	}
	q = q.OrderBy("s.name")

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
	q = q.Where(sq.Eq{"s.id": IDs})
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
		fmt.Println("INSERTING:")
		fmt.Println(m.SiteInfo.Geometry.EWKT(8))
		fmt.Println(m.SiteInfo.Geometry.Coordinates)
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
	q = q.Where("s.id = ?", siteID)
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

// func UpdateSite(db *pgxpool.Pool, s *Site) (*Site, error) {
// 	var id uuid.UUID
// 	if err := pgxscan.Get(
// 		context.Background(), db, &id,
// 		"UPDATE usgs_site SET name=$2, geometry=$3, elevation=$4, horizontal_datum_id=$5, vertical_datum_id=$6, huc=$7, state_id=$8, update_date=CURRENT_TIMESTAMP WHERE usgs_id = $1 RETURNING id",
// 		s.UsgsId, s.Name, s.Geometry.EWKT(), s.Elevation, s.HorizontalDatumId, s.VerticallDatumId, s.Huc, s.StateID,
// 	); err != nil {
// 		return nil, err
// 	}
// 	return GetSiteByID(db, &id)
// }

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
