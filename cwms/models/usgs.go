package models

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Site struct {
	ID                string     `json:"id"`
	UsgsId            string     `json:"usgs_id" db:"usgs_id"`
	Name              string     `json:"name"`
	Geometry          Geometry   `json:"geometry"`
	Elevation         *float32   `json:"elevation"`
	HorizontalDatumId int        `json:"horizontan_datum_id" db:"horizontal_datum_id"`
	VerticallDatumId  int        `json:"vertical_datum_id" db:"vertical_datum_id"`
	VerticalDatum     string     `json:"vertical_datum" db:"vertical_datum"`
	Huc               *string    `json:"huc"`
	StateID           *int       `json:"state_id" db:"state_id"`
	State             string     `json:"state"`
	StateAbbrev       string     `json:"state_abbrev" db:"state_abbrev"`
	CreateDate        time.Time  `json:"create_date" db:"create_date"`
	UpdateDate        *time.Time `json:"update_date" db:"update_date"`
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
					s.state_id,
					st.name as state,
					st.stusps as state_abbrev,
					s.create_date,
					s.update_date`,
	).From("usgs_site s").Join("vertical_datum v on v.id=s.vertical_datum_id").Join("tiger_data.state_all st on st.gid=s.state_id")

	if sf != nil {
		// Filter by StateID
		if sf.StateAbbrev != nil {
			q = q.Where("st.stusps = ?", strings.ToUpper(*sf.StateAbbrev))
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
