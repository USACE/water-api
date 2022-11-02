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
	VisualizationID *uuid.UUID `json:"visualization_id,omitempty" db:"chart_id"`
	Slug            string     `json:"slug,omitempty" db:"slug"`
	Variable        string     `json:"variable,omitempty" db:"variable"`
	Key             string     `json:"key,omitempty"`
	Datatype        string     `json:"datatype,omitempty" query:"datatype"`
	LatestTime      *time.Time `json:"latest_time,omitempty"`
	LatestValue     *float64   `json:"latest_value,omitempty"`
	Provider        string     `json:"provider,omitempty" db:"provider"`
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
							FROM visualization v
							LEFT JOIN "location" l ON l.id = v.location_id
							JOIN provider p ON p.id = v.provider_id`

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

	var getVisualizationSQL = `
								WITH timeseries_providers AS (
									-- make sure we get the provider attached
									-- to the timeseries instead of the viz/chart
									SELECT t.id, p.slug FROM provider p 
									JOIN datasource d ON d.provider_id = p.id
									JOIN timeseries t ON t.datasource_id = d.id 
								)
								SELECT
								l.slug AS location_slug,
								c.name AS name,
								c.slug AS slug,
								c.type_id,
								p."name" AS provider_name, 
								p.slug AS provider_slug,
								COALESCE(json_agg(json_build_object(
									'variable', cvm.variable,
									'datatype', dt.slug,
									'provider', tp.slug,
									'key', t.datasource_key,
									'latest_time', t.latest_time,
									'latest_value', t.latest_value
								)), '[]') AS mapping
								FROM chart c
								LEFT JOIN "location" l ON l.id = c.location_id
								JOIN provider p ON p.id = v.provider_id
								LEFT JOIN chart_variable_mapping cvm ON cvm.chart_id = c.id
								LEFT JOIN timeseries t ON t.id = cvm.timeseries_id
								LEFT JOIN datasource d ON d.id = t.datasource_id 
								LEFT JOIN datatype dt ON dt.id = d.datatype_id
								LEFT JOIN timeseries_providers tp ON tp.id = t.id
								WHERE lower(c.slug) = lower($1)
								GROUP BY l.slug, c.slug, c.name, c.type_id, p.name, p.slug
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
							c.name AS name,
							c.slug AS slug,
							c.type_id,
							p."name" AS provider_name, 
							p.slug AS provider_slug,
							COALESCE(json_agg(json_build_object(
								'variable', cvm.variable,
								'key', t.datasource_key,
								'latest_time', t.latest_time,
								'latest_value', t.latest_value
							)), '[]') AS mapping
							FROM chart c
							JOIN "location" l ON l.id = c.location_id
							JOIN provider p ON p.id = l.office_id 
							LEFT JOIN chart_variable_mapping cvm ON cvm.chart_id = c.id
							LEFT JOIN timeseries t ON t.id = cvm.timeseries_id
							WHERE lower(l.slug) = lower($1)
							AND type_id = $2
							GROUP BY l.slug, c.slug, c.name, c.type_id, p.name, p.slug
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
			`INSERT INTO chart (slug, name, type_id, provider_id)
			VALUES($1, $2, $3, (SELECT id from provider where lower(slug) = lower($4)))
			RETURNING slug`, slug, v.Name, v.TypeID, v.ProviderSlug,
		); err != nil {
			return nil, err
		}
	} else {
		if err := pgxscan.Get(
			context.Background(), db, &vSlug,
			`INSERT INTO chart (location_id, slug, name, type_id, provider_id)
			VALUES((SELECT l.id FROM location l WHERE lower(l.slug) = lower($1)), $2, $3, $4, (SELECT id from provider where lower(slug) = lower($5)))
			RETURNING slug`, v.LocationSlug, slug, v.Name, v.TypeID, v.ProviderSlug,
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
			`INSERT into chart_variable_mapping(chart_id, variable, timeseries_id) 
			VALUES
			(
				(SELECT id from chart WHERE lower(slug) = lower($1)), 
				$2,
				(SELECT t.id from timeseries t
					JOIN datasource d ON d.id = t.datasource_id 
					JOIN datatype dt ON dt.id = d.datatype_id
					JOIN provider p ON p.id = d.provider_id 
					WHERE lower(datasource_key) = lower($3)
					AND lower(dt.slug) = lower($4)
					AND lower(p.slug) = lower($5)
				)
			)
			ON CONFLICT ON CONSTRAINT chart_id_unique_variable
			DO UPDATE SET
			variable = $2
			WHERE chart_variable_mapping.chart_id = (SELECT id from chart WHERE lower(slug) = lower($1))
			RETURNING chart_id`,
			*visualizationSlug, v.Variable, v.Key, v.Datatype, v.Provider,
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

func DeleteVisualization(db *pgxpool.Pool, visualizationSlug *string) error {
	if _, err := db.Exec(context.Background(), `DELETE FROM chart WHERE slug=$1`, visualizationSlug); err != nil {
		return err
	}
	return nil
}
