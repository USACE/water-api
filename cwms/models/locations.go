package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

func (g Geometry) EWKT() string {
	return fmt.Sprintf("SRID=4326;POINT(%f %f)", g.Coordinates[0], g.Coordinates[1])
}

type Location struct {
	ID         uuid.UUID `json:"id"`
	OfficeID   uuid.UUID `json:"office_id"`
	StateID    *int      `json:"state_id"`
	Name       string    `json:"name"`
	PublicName *string   `json:"public_name"`
	Slug       string    `json:"slug"`
	KindID     uuid.UUID `json:"kind_id"`
	Kind       string    `json:"kind"`
	Geometry   Geometry  `json:"geometry"`
}

type LocationFilter struct {
	KindID   *uuid.UUID `json:"kind_id" query:"kind_id"`
	OfficeID *uuid.UUID `json:"office_id" query:"office_id"`
	StateID  *int       `json:"state_id" query:"state_id"`
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

func ListLocationsQuery(f *LocationFilter) (sq.SelectBuilder, error) {

	q := sq.Select(`a.id,
		            a.office_id,
					a.state_id,
		            a.name,
		            a.public_name,
		            a.slug,
		            ST_AsGeoJSON(a.geometry)::json AS geometry,
		            k.id                           AS kind_id,
		            k.name                         AS kind`,
	).From("location a")

	// Base string for JOIN of location_kind table
	j1 := "location_kind k ON k.id = a.kind_id"

	if f != nil {
		// Filter by KindID
		if f.KindID != nil {
			// limit join table kind if kind_id provided
			q = q.Join(fmt.Sprintf("%s AND k.id = ?", j1), f.KindID).Where("k.id = ?", f.KindID)
		} else {
			// always join location_kind
			q = q.Join(j1)
		}
		// Filter by OfficeID
		if f.OfficeID != nil {
			q = q.Where("a.office_id = ?", f.OfficeID)
		}
		// Filter by StateID
		if f.StateID != nil {
			q = q.Where("a.state_id = ?", f.StateID)
		}
	} else {
		// always join location_kind
		q = q.Join(j1)
	}

	// Unfiltered
	return q.PlaceholderFormat(sq.Dollar), nil
}

func ListLocations(db *pgxpool.Pool, f *LocationFilter) ([]Location, error) {
	q, err := ListLocationsQuery(f)
	if err != nil {
		return make([]Location, 0), err
	}
	sql, args, err := q.ToSql()
	if err != nil {
		return make([]Location, 0), err
	}
	ll := make([]Location, 0)
	if err := pgxscan.Select(context.Background(), db, &ll, sql, args...); err != nil {
		return make([]Location, 0), err
	}
	return ll, nil
}

func ListLocationsForIDs(db *pgxpool.Pool, IDs []uuid.UUID) ([]Location, error) {
	// Base Locations Query
	q, err := ListLocationsQuery(nil)
	if err != nil {
		return make([]Location, 0), err
	}
	// Where ID In (...)
	q = q.Where(sq.Eq{"a.id": IDs})
	sql, args, err := q.ToSql()
	if err != nil {
		return make([]Location, 0), err
	}
	ll := make([]Location, 0)
	if err := pgxscan.Select(context.Background(), db, &ll, sql, args...); err != nil {
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
			m.OfficeID, m.Name, m.PublicName, m.Slug, m.Geometry.EWKT(), m.KindID,
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

func GetLocationByID(db *pgxpool.Pool, locationID *uuid.UUID) (*Location, error) {
	// Base Locations Query
	q, err := ListLocationsQuery(nil)
	if err != nil {
		return nil, err
	}
	// Where ID In (...)
	q = q.Where("a.id = ?", locationID)
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	var l Location
	if err := pgxscan.Get(context.Background(), db, &l, sql, args...); err != nil {
		return nil, err
	}
	return &l, nil
}

func GetLocationBySlug(db *pgxpool.Pool, locationSlug *string) (*Location, error) {
	// Base Locations Query
	q, err := ListLocationsQuery(nil)
	if err != nil {
		return nil, err
	}
	// Where slug =
	q = q.Where("a.slug = ?", locationSlug)
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	var l Location
	if err := pgxscan.Get(context.Background(), db, &l, sql, args...); err != nil {
		return nil, err
	}
	return &l, nil
}

func UpdateLocation(db *pgxpool.Pool, l *Location) (*Location, error) {
	var id uuid.UUID
	if err := pgxscan.Get(
		context.Background(), db, &id,
		"UPDATE location SET update_date=CURRENT_TIMESTAMP, office_id=$2, name=$3, public_name=$4, geometry=$5, kind_id=$6 WHERE id = $1 RETURNING id",
		l.ID, l.OfficeID, l.Name, l.PublicName, l.Geometry.EWKT(), l.KindID,
	); err != nil {
		return nil, err
	}
	return GetLocationByID(db, &id)
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
