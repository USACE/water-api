package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/USACE/water-api/api/app"
	"github.com/USACE/water-api/api/helpers"
	"github.com/USACE/water-api/api/timeseries"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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

// Extract struct for the return of measurements for a watershed
type Extract struct {
	SiteNumber  string      `json:"site_number"`
	Code        string      `json:"code"`
	SiteID      uuid.UUID   `json:"site_id"`
	Name        string      `json:"name"`
	ParameterID uuid.UUID   `json:"paramter_id"`
	Times       []time.Time `json:"times"`
	Values      []float64   `json:"values"`
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
func GetWatershed(db *pgxpool.Pool, watershedSlug *string) (*Watershed, error) {
	var w Watershed
	if err := pgxscan.Get(
		context.Background(), db, &w, WatershedSQL+` FROM v_watershed w WHERE w.slug = $1`, watershedSlug,
	); err != nil {
		return nil, err
	}
	return &w, nil
}

// GetWatershedGeometry returns a single watershed geometry (as GeoJSON) using slug
func GetWatershedGeometry(db *pgxpool.Pool, watershedSlug *string) ([]byte, error) {
	var j []byte
	if err := pgxscan.Get(
		context.Background(), db, &j, ` 
			select json_build_object(
				'type', 'FeatureCollection',
				'features', json_agg(ST_AsGeoJSON(t.*, 'geom', 6, false)::json)
			)
			from (SELECT w.id, w.name, ST_ForcePolygonCCW(ST_Transform(w.geometry,4326)) 
				  FROM a2w_cwms.watershed w WHERE slug = $1
			) as t(id, name, geom)
			`, watershedSlug,
	); err != nil {
		return nil, err
	}
	return j, nil
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
	return GetWatershed(db, &wNew.Slug)
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
		`UPDATE watershed SET geometry = ST_Transform(ST_GeomFromGeoJSON($1),4326) WHERE id = $2 RETURNING slug`,
		qs,
		id,
	)
	if err != nil {
		return nil, err
	}
	var slug string
	if err = pgxscan.ScanOne(&slug, rows); err != nil {
		return nil, err
	}
	// Commit and then get the updated watersed to return
	if err = tx.Commit(context.Background()); err != nil {
		return nil, err
	}
	ws, err := GetWatershed(db, &slug)
	if err != nil {
		return nil, err
	}
	return ws, err
}

// UpdateWatershed updates a watershed
func UpdateWatershed(db *pgxpool.Pool, w *Watershed) (*Watershed, error) {
	var wSlug string
	if err := pgxscan.Get(context.Background(), db, &wSlug, `UPDATE watershed SET name=$1, office_id=$3 WHERE slug=$2 RETURNING slug`, &w.Name, &w.Slug, &w.OfficeID); err != nil {
		return nil, err
	}
	fmt.Println(wSlug)
	return GetWatershed(db, &wSlug)
}

// DeleteWatershed deletes a watershed by id
func DeleteWatershed(db *pgxpool.Pool, watershedSlug *string) error {
	if _, err := db.Exec(context.Background(), `UPDATE watershed SET deleted=true WHERE slug=$1`, watershedSlug); err != nil {
		return err
	}
	return nil
}

func UndeleteWatershed(db *pgxpool.Pool, watershedSlug *string) (*Watershed, error) {
	var wSlug string
	if err := pgxscan.Get(
		context.Background(), db, &wSlug, `UPDATE watershed SET deleted=false WHERE slug=$1 RETURNING slug`, watershedSlug,
	); err != nil {
		return nil, err
	}
	return GetWatershed(db, &wSlug)
}

// UploadWatersheds
func UploadWatersheds(db *pgxpool.Pool, wid uuid.UUID, file *multipart.FileHeader) (map[string]string, error) {
	// Begin DB transaction getting the watershed slug from the uuid
	tx, err := db.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err != nil {
		return nil, err
	}
	var slug string
	if err = pgxscan.Get(context.Background(), db, &slug, `SELECT slug FROM watershed WHERE id = $1`, wid); err != nil {
		return nil, err
	}
	// Upload file and add record to database
	fn := file.Filename
	fs := file.Size
	key := aws.String("water/watersheds/" + slug + "/" + fn)
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buffer := make([]byte, fs)
	_, err = f.Read(buffer)
	if err != nil {
		return nil, err
	}
	cfg, err := app.GetConfig()
	if err != nil {
		return nil, err
	}
	bucket := aws.String(cfg.AWSS3Bucket)
	s3Config := app.AWSConfig(*cfg)
	newSession, _ := session.NewSession(&s3Config)
	s3Client := s3.New(newSession)
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(buffer),
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		return nil, err
	}
	// Add record to database
	var id uuid.UUID
	if err = db.QueryRow(context.Background(),
		`INSERT INTO watershed_shapefile_uploads (watershed_id, file, file_size)
		VALUES($1, $2, $3)
		RETURNING id`,
		wid, key, fs,
	).Scan(&id); err != nil {
		return nil, err
	}

	result := map[string]string{}
	result["id"] = id.String()
	result["Filename"] = fn
	result["FileSize"] = fmt.Sprint(fs) + " bytes"
	result["Key"] = *key
	result["Bucket"] = *bucket

	return result, nil
}

// WatershedExtract
func WatershedExtract(db *pgxpool.Pool, slug string, tw *timeseries.TimeWindow) ([]Extract, error) {
	ext := make([]Extract, 0)
	rows, err := db.Query(context.Background(),
		`SELECT r1.site_number, r1."name", r1.code, r1.site_id, r1.parameter_id, array_agg(r1."time") AS times, array_agg(r1.value) AS values
        FROM
        (
            SELECT us.site_number, us."name", up.code, usp.site_id, usp.parameter_id, um."time", um.value 
            FROM watershed_usgs_sites wus
            JOIN
            (
            SELECT "time", value, usgs_site_parameters_id from usgs_measurements
            WHERE "time" >= $2::timestamptz AND "time" <= $3::timestamptz
            ORDER BY usgs_site_parameters_id, "time" desc
            ) AS um on um.usgs_site_parameters_id = wus.usgs_site_parameter_id
        JOIN
        usgs_site_parameters usp 
        ON usgs_site_parameters_id = usp.id
        JOIN usgs_site us 
        ON us.id = usp.site_id
        JOIN usgs_parameter up
        ON up.id = usp.parameter_id
        WHERE wus.watershed_id = (SELECT id FROM watershed w WHERE slug = $1)
        ORDER BY us.site_number, up.code, um."time"
        ) r1
        GROUP BY r1.site_number, r1."name", r1.code, r1.site_id, r1.parameter_id`,
		slug, tw.After, tw.Before,
	)
	if err != nil {
		return nil, err
	}
	if err = pgxscan.ScanAll(&ext, rows); err != nil {
		return nil, err
	}

	return ext, nil
}
