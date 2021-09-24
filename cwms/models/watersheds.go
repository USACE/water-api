package models

import (
	"context"
	"encoding/json"

	"github.com/USACE/water-api/helpers"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/paulmach/orb/geojson"
	jsonpointer "github.com/xeipuuv/gojsonpointer"
)

// Watershed is a watershed struct
type Watershed struct {
	ID           uuid.UUID `json:"id" db:"id"`
	OfficeSymbol *string   `json:"office_symbol" db:"office_symbol"`
	OfficeID     uuid.UUID `json:"office_id" db:"office_id"`
	Slug         string    `json:"slug"`
	Name         string    `json:"name"`
	Bbox         []float64 `json:"bbox" db:"bbox"`
}

// WatershedSQL includes common fields selected to build a watershed
const WatershedSQL = `SELECT w.id,
                             w.office_symbol,
							 w.office_id,
                             w.slug,
                             w.name,
	                         ARRAY[
								 ST_XMin(w.geometry),
								 ST_Ymin(w.geometry),
								 ST_XMax(w.geometry),
								 ST_YMax(w.geometry)
							 ] AS bbox`

// ListWatersheds returns an array of watersheds
func ListWatersheds(db *pgxpool.Pool) ([]Watershed, error) {
	ww := make([]Watershed, 0)
	if err := pgxscan.Select(context.Background(), db, &ww, WatershedSQL+" FROM v_watershed w order by w.office_symbol, w.name"); err != nil {
		return make([]Watershed, 0), nil
	}
	return ww, nil
}

// GetWatershed returns a single watershed using slug
func GetWatershed(db *pgxpool.Pool, watershedID *uuid.UUID) (*Watershed, error) {
	var w Watershed
	if err := pgxscan.Get(
		context.Background(), db, &w, WatershedSQL+` FROM v_watershed w WHERE w.id = $1`, watershedID,
	); err != nil {
		return nil, err
	}
	return &w, nil
}

// CreateWatershed creates a new watershed
func CreateWatershed(db *pgxpool.Pool, w *Watershed) (*Watershed, error) {
	slug, err := helpers.NextUniqueSlug(db, "watershed", "slug", w.Name, "", "")
	if err != nil {
		return nil, err
	}
	var wNew Watershed
	if err := pgxscan.Get(
		context.Background(), db, &wNew,
		`INSERT INTO watershed (name, slug, office_id) VALUES ($1,$2, $3) RETURNING id, name, slug, office_id`, &w.Name, slug, &w.OfficeID,
	); err != nil {
		return nil, err
	}
	return GetWatershed(db, &wNew.ID)
	//return &wNew, nil
}

// UpdateWatershedGeometry
func UpdateWatershedGeometry(db *pgxpool.Pool, id *uuid.UUID, wf *geojson.Feature) (*Watershed, error) {
	// Database transaction
	tx, err := db.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err != nil {
		return nil, err
	}
	// Extract the GeoJSON from the Feature
	j, err := wf.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var json_doc map[string]interface{}
	if err = json.Unmarshal(j, &json_doc); err != nil {
		return nil, err
	}
	pointer_string := "/geometry/coordinates"
	pointer, _ := jsonpointer.NewJsonPointer(pointer_string)
	geo, _, err := pointer.Get(json_doc)
	if err != nil {
		return nil, err
	}
	geo_coords, err := json.Marshal(geo)
	if err != nil {
		return nil, err
	}
	// Build the query and UPDATE geometry
	geo_type := wf.Geometry.GeoJSONType()
	qs := "{\"type\":\"" + geo_type + "\","
	qs += "\"coordinates\":" + string(geo_coords)
	qs += "}"

	rows, err := tx.Query(
		context.Background(),
		`UPDATE watershed SET geometry = ST_Transform(ST_GeomFromGeoJSON($1),4326) WHERE id = $2 RETURNING id`,
		qs,
		id,
	)
	if err != nil {
		return nil, err
	}
	var rid uuid.UUID
	if err = pgxscan.ScanOne(&rid, rows); err != nil {
		return nil, err
	}
	// Commit and then get the updated watersed to return
	if err = tx.Commit(context.Background()); err != nil {
		return nil, err
	}
	ws, err := GetWatershed(db, &rid)
	if err != nil {
		return nil, err
	}
	return ws, err
}

// UpdateWatershed updates a watershed
func UpdateWatershed(db *pgxpool.Pool, w *Watershed) (*Watershed, error) {
	var wID uuid.UUID
	if err := pgxscan.Get(context.Background(), db, &wID, `UPDATE watershed SET name=$1, office_id=$3 WHERE id=$2 RETURNING id`, &w.Name, &w.ID, &w.OfficeID); err != nil {
		return nil, err
	}
	return GetWatershed(db, &wID)
}

// DeleteWatershed deletes a watershed by slug
func DeleteWatershed(db *pgxpool.Pool, watershedID *uuid.UUID) error {
	if _, err := db.Exec(context.Background(), `UPDATE watershed SET deleted=true WHERE id=$1`, watershedID); err != nil {
		return err
	}
	return nil
}

func UndeleteWatershed(db *pgxpool.Pool, watershedID *uuid.UUID) (*Watershed, error) {
	var wID uuid.UUID
	if err := pgxscan.Get(
		context.Background(), db, &wID, `UPDATE watershed SET deleted=false WHERE id=$1 RETURNING id`, watershedID,
	); err != nil {
		return nil, err
	}
	return GetWatershed(db, &wID)
}
