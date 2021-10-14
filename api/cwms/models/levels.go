package models

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Level struct
type (
	Level struct {
		ID   uuid.UUID `json:"id" db:"id"`
		Slug string    `json:"slug" db:"slug"`
		Name string    `json:"name" db:"name"`
	}
	LevelValues struct {
		Time  time.Time `json:"time"`
		Value float64   `json:"value"`
	}
	Levels struct {
		KindID uuid.UUID     `json:"kind_id"`
		Levels []LevelValues `json:"levels"`
	}
)

// ListLevelKind
func ListLevelKind(db *pgxpool.Pool) ([]Level, error) {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	tx.Rollback(context.Background())

	lvl := []Level{}
	if err = pgxscan.Select(context.Background(),
		db,
		&lvl,
		`SELECT id, slug, name FROM level_kind`,
	); err != nil {
		return nil, err
	}
	return lvl, nil
}

// CreateLevelKind
func CreateLevelKind(db *pgxpool.Pool, slug_name string, name string) (*Level, error) {
	// Begin transaction
	tx, err := db.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err != nil {
		return nil, err
	}
	// Insert into returning row
	var lvl Level
	if err = pgxscan.Get(context.Background(),
		db,
		&lvl,
		`INSERT INTO level_kind(slug, name) VALUES ($1, $2) RETURNING id, slug, name`,
		slug_name, name,
	); err != nil {
		return nil, err
	}
	// Commit the transaction or return err
	if err = tx.Commit(context.Background()); err != nil {
		return nil, err
	}
	return &lvl, nil
}

// DeleteLevelKind
func DeleteLevelKind(db *pgxpool.Pool, ls string) (int64, error) {
	// Delete entry by slug
	res, err := db.Exec(context.Background(),
		`DELETE FROM level_kind WHERE id IN (SELECT id FROM level_kind WHERE slug = $1) RETURNING name`,
		ls,
	)
	if err != nil {
		return 0, err
	}
	count := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CreateLocationLevels
func CreateLocationLevels(db *pgxpool.Pool, loc_id uuid.UUID, p []Levels) error {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	for idx := range p {
		k := p[idx]
		kind_id := k.KindID
		lvls := k.Levels
		for ndx := range lvls {
			lvl := lvls[ndx]
			if _, err = db.Exec(context.Background(),
				`INSERT INTO level_value (level_id, julian_date, value)
				VALUES (
					(SELECT id FROM level WHERE location_id = $1 AND level_kind_id = $2),
					(SELECT EXTRACT(epoch FROM $3::timestamptz)),
					$4
				)`,
				loc_id, kind_id, lvl.Time, lvl.Value,
			); err != nil {
				return err
			}
		}
	}
	if err = tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}

// UpdateLocationLevels
func UpdateLocationLevels(db *pgxpool.Pool, loc_id uuid.UUID, p []Levels) error {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	for idx := range p {
		k := p[idx]
		kind_id := k.KindID
		lvls := k.Levels
		for ndx := range lvls {
			lvl := lvls[ndx]
			if _, err = db.Exec(context.Background(),
				`UPDATE level_value
				SET value = $1
				WHERE level_id = (SELECT id FROM level WHERE location_id = $2 AND
								  level_kind_id = $3) AND
				julian_date = (SELECT EXTRACT(epoch FROM $4::timestamptz));`,
				lvl.Value, loc_id, kind_id, lvl.Time,
			); err != nil {
				return err
			}
		}
	}
	if err = tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil
}

// ListLevelValues
func ListLevelValues(db *pgxpool.Pool, slug string, kind string) ([]LevelValues, error) {
	var lvls []LevelValues
	rows, err := db.Query(
		context.Background(),
		`SELECT to_timestamp(julian_date) AS time, value
			FROM level_value
			WHERE level_id = (SELECT id
				FROM level
				WHERE location_id = (SELECT id FROM location WHERE slug = $1) AND
				level_kind_id = (SELECT id FROM level_kind WHERE slug = $2)
				)`,
		slug, kind,
	)
	if err != nil {
		return nil, err
	}
	if err = pgxscan.ScanAll(&lvls, rows); err != nil {
		fmt.Println(err.Error())
	}
	return lvls, nil
}
