package models

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/USACE/water-api/api/helpers"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type VisualizationMapping struct {
	VisualizationID *uuid.UUID `json:"visualization_id,omitempty" db:"visualization_id"`
	Slug            string     `json:"slug,omitempty" db:"slug"`
	Variable        string     `json:"variable,omitempty" db:"variable"`
	Key             string     `json:"key,omitempty"`
	DatasourceType  string     `json:"datasource_type,omitempty" query:"datasource_type"`
	LatestTime      *time.Time `json:"latest_time,omitempty"`
	LatestValue     *float64   `json:"latest_value,omitempty"`
}

type Visualization struct {
	LocationSlug *string   `json:"location_slug" db:"location_slug"`
	Name         string    `json:"name" db:"name"`
	Slug         string    `json:"slug" db:"slug"`
	TypeID       uuid.UUID `json:"type_id" db:"type_id"`
	ProviderName string    `json:"provider_name" db:"provider_name"`
	ProviderSlug string    `json:"provider_slug" db:"provider_slug"`

	Mapping *[]VisualizationMapping `json:"mapping,omitempty"`
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

var listVisualizationsSQL = `SELECT 
							l.slug AS location_slug, 
							v.name, v.slug, v.type_id, 
							p."name" AS provider_name, 
							p.slug AS provider_slug
							FROM a2w_cwms.visualization v 
							JOIN a2w_cwms."location" l ON l.id = v.location_id
							JOIN a2w_cwms.provider p ON p.id = l.office_id`

func ListVisualizations(db *pgxpool.Pool) ([]Visualization, error) {
	vv := make([]Visualization, 0)
	if err := pgxscan.Select(
		context.Background(), db, &vv, listVisualizationsSQL+" ORDER BY p.slug, v.slug",
	); err != nil {
		return make([]Visualization, 0), err
	}
	return vv, nil
}

// GetVisualization returns a single Visualization
func GetVisualization(db *pgxpool.Pool, visualizationSlug *string) (*Visualization, error) {

	var getVisualizationSQL = `SELECT
							l.slug AS location_slug,
							v.name AS name,
							v.slug AS slug,
							v.type_id,
							p."name" AS provider_name, 
							p.slug AS provider_slug,
							COALESCE(json_agg(json_build_object(
								'variable', vvm.variable,
								'key', t.datasource_key,
								'latest_time', t.latest_time,
								'latest_value', t.latest_value
							)), '[]') AS mapping
							FROM a2w_cwms.visualization v
							JOIN a2w_cwms."location" l ON l.id = v.location_id
							JOIN provider p ON p.id = l.office_id 
							LEFT JOIN a2w_cwms.visualization_variable_mapping vvm ON vvm.visualization_id = v.id
							LEFT JOIN a2w_cwms.timeseries t ON t.id = vvm.timeseries_id
							WHERE lower(v.slug) = lower($1)
							GROUP BY l.slug, v.slug, v.name, v.type_id, p.name, p.slug
							LIMIT 1`

	var v Visualization
	//if err := pgxscan.Get(context.Background(), db, &v, listVisualizationsSQL+" WHERE lower(v.slug) = lower($1)", visualizationSlug); err != nil {
	if err := pgxscan.Get(context.Background(), db, &v, getVisualizationSQL, visualizationSlug); err != nil {
		return nil, err
	}
	return &v, nil
}

func GetVisualizationByLocation(db *pgxpool.Pool, locationSlug *string, visualizationTypeId *uuid.UUID) (*Visualization, error) {

	var getVisualizationSQL = `SELECT
							l.slug AS location_slug,
							v.name AS name,
							v.slug AS slug,
							v.type_id,
							p."name" AS provider_name, 
							p.slug AS provider_slug,
							COALESCE(json_agg(json_build_object(
								'variable', vvm.variable,
								'key', t.datasource_key,
								'latest_time', t.latest_time,
								'latest_value', t.latest_value
							)), '[]') AS mapping
							FROM a2w_cwms.visualization v
							JOIN a2w_cwms."location" l ON l.id = v.location_id
							JOIN provider p ON p.id = l.office_id 
							LEFT JOIN a2w_cwms.visualization_variable_mapping vvm ON vvm.visualization_id = v.id
							LEFT JOIN a2w_cwms.timeseries t ON t.id = vvm.timeseries_id
							WHERE lower(l.slug) = lower($1)
							AND type_id = $2
							GROUP BY l.slug, v.slug, v.name, v.type_id, p.name, p.slug
							LIMIT 1`

	var v Visualization
	if err := pgxscan.Get(context.Background(), db, &v, getVisualizationSQL, locationSlug, visualizationTypeId); err != nil {
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

	var vSlug string
	if v.LocationSlug == nil {
		if err := pgxscan.Get(
			context.Background(), db, &vSlug,
			`INSERT INTO visualization (slug, name, type_id)
			VALUES($1, $2, $3)
			RETURNING slug`, slug, v.Name, v.TypeID,
		); err != nil {
			return nil, err
		}
	} else {
		if err := pgxscan.Get(
			context.Background(), db, &vSlug,
			`INSERT INTO visualization (location_id, slug, name, type_id)
			VALUES((SELECT l.id FROM location l WHERE lower(l.slug) = lower($1)), $2, $3, $4)
			RETURNING slug`, v.LocationSlug, slug, v.Name, v.TypeID,
		); err != nil {
			return nil, err
		}
	}

	return GetVisualization(db, &vSlug)

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
			)
			ON CONFLICT ON CONSTRAINT visualization_id_unique_variable
			DO UPDATE SET
			variable = $2
			WHERE visualization_variable_mapping.visualization_id = (SELECT id from visualization WHERE lower(slug) = lower($1))
			RETURNING visualization_id`,
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
