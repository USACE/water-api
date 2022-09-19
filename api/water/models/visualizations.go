package models

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/USACE/water-api/api/helpers"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type VisualizationMapping struct {
	VisualizationID uuid.UUID `json:"visualization_id" db:"visualization_id"`
	Slug            string    `json:"slug" db:"slug"`
	Variable        string    `json:"variable" db:"variable"`
	Key             string    `json:"key"`
	DatasourceType  string    `json:"datasource_type" query:"datasource_type"`
}

type Visualization struct {
	LocationSlug string    `json:"location_slug" db:"location_slug"`
	Name         string    `json:"name" db:"name"`
	Slug         string    `json:"slug" db:"slug"`
	TypeID       uuid.UUID `json:"type_id" db:"type_id"`
}

type VisualizationMappingCollection struct {
	Items []VisualizationMapping `json:"items"`
}

func (c *VisualizationMappingCollection) UnmarshalJSON(b []byte) error {
	switch helpers.JSONType(b) {
	case "ARRAY":
		return json.Unmarshal(b, &c.Items)
	case "OBJECT":
		c.Items = make([]VisualizationMapping, 1)
		return json.Unmarshal(b, &c.Items[0])
	default:
		return errors.New("payload not recognized as JSON array or object")
	}
}

var listVisualizationsSQL = `SELECT l.slug AS location_slug, v.name, v.slug, v.type_id
						FROM a2w_cwms.visualization v 
						JOIN a2w_cwms."location" l ON l.id = v.location_id`

func ListVisualizations(db *pgxpool.Pool) ([]Visualization, error) {
	vv := make([]Visualization, 0)
	if err := pgxscan.Select(
		context.Background(), db, &vv, listVisualizationsSQL+" ORDER BY v.slug",
	); err != nil {
		return make([]Visualization, 0), err
	}
	return vv, nil
}

// GetVisualization returns a single Visualization
func GetVisualization(db *pgxpool.Pool, visualizationID *uuid.UUID) (*Visualization, error) {
	var v Visualization
	if err := pgxscan.Get(context.Background(), db, &v, listVisualizationsSQL+" WHERE v.id = $1", visualizationID); err != nil {
		return nil, err
	}
	return &v, nil
}

// CreateVisualization creates a single Visualization
func CreateVisualization(db *pgxpool.Pool, v *Visualization) (*Visualization, error) {

	// Insert Into Database Using New Slug
	slug, err := helpers.NextUniqueSlug(db, "visualization", "slug", v.Name, "", "")
	if err != nil {
		return nil, err
	}

	var vID uuid.UUID
	if err := pgxscan.Get(
		context.Background(), db, &vID,
		`INSERT INTO visualization (location_id, slug, name, type_id)
		VALUES((SELECT l.id FROM location l WHERE lower(l.slug) = lower($1)), $2, $3, $4)
		RETURNING id`, v.LocationSlug, slug, v.Name, v.TypeID,
	); err != nil {
		return nil, err
	}
	return GetVisualization(db, &vID)

}

func CreateOrUpdateVisualizationMapping(db *pgxpool.Pool, c VisualizationMappingCollection, visualizationSlug *string) ([]VisualizationMapping, error) {

	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]VisualizationMapping, 0), err
	}
	defer tx.Rollback(context.Background())

	updatedIDs := make([]uuid.UUID, 0)

	for _, v := range c.Items {

		rows, err := tx.Query(
			context.Background(),
			`INSERT into visualization_variable_mapping(visualization_id, variable, timeseries_id) 
			VALUES
			(
				(SELECT id from visualization WHERE lower(slug) = lower($1)), 
				$2,
				(SELECT t.id from timeseries t
					JOIN a2w_cwms.datasource d ON d.id = t.datasource_id 
					JOIN a2w_cwms.datasource_type dt ON dt.id = d.datasource_type_id 
					WHERE lower(datasource_key) = lower($3)
					AND lower(dt.slug) = lower($4))
			) RETURNING visualization_id`,
			*visualizationSlug, v.Variable, v.Key, v.DatasourceType,
		)
		if err != nil {
			return make([]VisualizationMapping, 0), err
		}
		var id uuid.UUID
		if err := pgxscan.ScanOne(&id, rows); err != nil {
			tx.Rollback(context.Background())
			return c.Items, err
		} else {
			updatedIDs = append(updatedIDs, id)
		}

	}
	tx.Commit(context.Background())

	if len(updatedIDs) == 0 {
		return nil, errors.New("no records updated")
	}

	return nil, nil

}
