package models

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

type LocationDetail struct {
	// Inherit all fields from Location
	// ID         uuid.UUID `json:"id"`
	// OfficeID   uuid.UUID `json:"office_id"`
	// StateID    *int      `json:"state_id"`
	// Name       string    `json:"name"`
	// PublicName *string   `json:"public_name"`
	// Slug       string    `json:"slug"`
	// KindID     uuid.UUID `json:"kind_id"`
	// Kind       string    `json:"kind"`
	// Geometry   Geometry  `json:"geometry"`
	Location
	Image string `json:"image"`
}

func GetLocationDetail(db *pgxpool.Pool, locSlug *string) (*LocationDetail, error) {
	var d LocationDetail
	if err := pgxscan.Get(
		context.Background(), db, &d,
		`SELECT a.id,
		        a.office_id,
				a.state_id,
		        a.name,
		        a.public_name,
		        a.slug,
		        ST_AsGeoJSON(a.geometry)::json  AS geometry,
		        k.id                            AS kind_id,
		        k.name                          AS kind,
				'http://localhost:3000/dam.jpg' AS image
		FROM location a
		JOIN location_kind k ON k.id = a.kind_id
		WHERE slug = $1`, locSlug,
	); err != nil {
		return nil, err
	}
	return &d, nil
}
