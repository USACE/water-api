package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type QueryAndParams struct {
	Query  string
	Params []interface{}
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

func (g Geometry) WKT() string {
	return fmt.Sprintf("POINT(%f %f)", g.Coordinates[0], g.Coordinates[1])
}

type Location struct {
	ID         uuid.UUID `json:"id"`
	OfficeID   uuid.UUID `json:"office_id"`
	Name       string    `json:"name"`
	PublicName *string   `json:"public_name"`
	Slug       string    `json:"slug"`
	KindID     uuid.UUID `json:"kind_id"`
	Kind       string    `json:"kind"`
	Geometry   Geometry  `json:"geometry"`
}

type LocationFilter struct {
	OfficeID *uuid.UUID `json:"office_id" query:"office_id"`
	KindID   *uuid.UUID `json:"kind_id" query:"kind_id"`
}

type LocationCollection struct {
	Items []Location `json:"items"`
}

func (c *LocationCollection) UnmarshalJSON(b []byte) error {
	switch JSONType(b) {
	case "ARRAY":
		return json.Unmarshal(b, &c.Items)
	case "OBJECT":
		c.Items = make([]Location, 1)
		return json.Unmarshal(b, &c.Items[0])
	default:
		return errors.New("payload not recognized as JSON array or object")
	}
}

var ListLocationsBaseSQL = `SELECT a.id,
                                   a.office_id,
                                   a.name,
                                   a.public_name,
                                   a.slug,
                                   ST_AsGeoJSON(a.geometry)::json AS geometry,
                                   k.id as kind_id,
                                   k.name AS kind
							FROM location a
							JOIN location_kind k on k.id = a.kind_id`

func ListLocationsSQL(f *LocationFilter) QueryAndParams {

	if f != nil {
		// Filter by KindID and OfficeID
		if f.KindID != nil && f.OfficeID != nil {
			return QueryAndParams{
				Query:  ListLocationsBaseSQL + ` AND k.id = $1 WHERE k.id = $1 AND a.office_id = $2`,
				Params: append([]interface{}{}, f.KindID, f.OfficeID),
			}
		}
		// Filter by Only OfficeID
		if f.OfficeID != nil {
			return QueryAndParams{
				Query:  ListLocationsBaseSQL + ` WHERE a.office_id = $1`,
				Params: append([]interface{}{}, f.OfficeID),
			}
		}
		// Filter by Only KindID
		if f.KindID != nil {
			return QueryAndParams{
				Query:  ListLocationsBaseSQL + ` WHERE k.id = $1`,
				Params: append([]interface{}{}, f.KindID),
			}
		}
	}

	// Unfiltered
	return QueryAndParams{Query: ListLocationsBaseSQL, Params: []interface{}{}}
}

func ListLocations(db *pgxpool.Pool, f *LocationFilter) ([]Location, error) {
	ll, q := make([]Location, 0), ListLocationsSQL(f)
	if err := pgxscan.Select(context.Background(), db, &ll, q.Query, q.Params...); err != nil {
		return make([]Location, 0), err
	}
	return ll, nil
}

func ListLocationsForIDs(db *pgxpool.Pool, IDs []uuid.UUID) ([]Location, error) {
	ll := make([]Location, 0)
	s := ListLocationsSQL(nil)
	if err := pgxscan.Select(context.Background(), db, &ll, s.Query+" WHERE a.id = ANY($1)", IDs); err != nil {
		return make([]Location, 0), err
	}
	return ll, nil
}

func CreateLocations(db *pgxpool.Pool, n LocationCollection) ([]Location, error) {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]Location, 0), err
	}
	defer tx.Rollback(context.Background())
	newIDs := make([]uuid.UUID, 0)
	for _, m := range n.Items {
		rows, err := tx.Query(
			context.Background(),
			`INSERT INTO location (office_id, name, public_name, slug, geometry, kind_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
			m.OfficeID, m.Name, m.PublicName, m.Slug, m.Geometry.WKT(), m.KindID,
		)
		if err != nil {
			tx.Rollback(context.Background())
			return make([]Location, 0), err
		}
		var id uuid.UUID
		if err := pgxscan.ScanOne(&id, rows); err != nil {
			tx.Rollback(context.Background())
			return make([]Location, 0), err
		} else {
			newIDs = append(newIDs, id)
		}
	}
	tx.Commit(context.Background())

	return ListLocationsForIDs(db, newIDs)
}

func GetLocation(db *pgxpool.Pool, locationID *uuid.UUID) (*Location, error) {
	var l Location
	if err := pgxscan.Get(context.Background(), db, &l, ListLocationsBaseSQL+" WHERE a.id = $1", locationID); err != nil {
		return nil, err
	}
	return &l, nil
}

func UpdateLocation(db *pgxpool.Pool, l *Location) (*Location, error) {
	var id uuid.UUID
	if err := pgxscan.Get(
		context.Background(), db, &id,
		"UPDATE location SET update_date=CURRENT_TIMESTAMP, office_id=$2, name=$3, public_name=$4, geometry=$5, kind_id=$6 WHERE id = $1 RETURNING id",
		l.ID, l.OfficeID, l.Name, l.PublicName, l.Geometry.WKT(), l.KindID,
	); err != nil {
		return nil, err
	}
	return GetLocation(db, &id)
}

func DeleteLocation(db *pgxpool.Pool, locationID *uuid.UUID) error {
	if _, err := db.Exec(context.Background(), `DELETE FROM location WHERE id = $1`, locationID); err != nil {
		return err
	}
	return nil
}

func ListProjects(db *pgxpool.Pool) ([]Location, error) {
	// Known UUID For location_kind = 'PROJECT'
	kindID, err := uuid.Parse("460ea73b-c65e-4fc8-907a-6e6fd2907a99")
	if err != nil {
		return make([]Location, 0), err
	}
	return ListLocations(db, &LocationFilter{KindID: &kindID})
}
