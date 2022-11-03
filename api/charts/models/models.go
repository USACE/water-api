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

type ChartMapping struct {
	ChartID     *uuid.UUID `json:"chart_id,omitempty" db:"chart_id"`
	Slug        string     `json:"slug,omitempty" db:"slug"`
	Variable    string     `json:"variable,omitempty" db:"variable"`
	Key         string     `json:"key,omitempty"`
	Datatype    string     `json:"datatype,omitempty" query:"datatype"`
	LatestTime  *time.Time `json:"latest_time,omitempty"`
	LatestValue *float64   `json:"latest_value,omitempty"`
	Provider    string     `json:"provider,omitempty" db:"provider"`
}

type Chart struct {
	LocationSlug *string   `json:"location_slug" db:"location_slug"`
	Name         string    `json:"name" db:"name"`
	Slug         string    `json:"slug" db:"slug"`
	TypeID       uuid.UUID `json:"type_id" db:"type_id"`
	ProviderName string    `json:"provider_name" db:"provider_name"`
	ProviderSlug string    `json:"provider_slug" db:"provider_slug"`

	Mapping *[]ChartMapping `json:"mapping,omitempty"`
}

type ChartMappingCollection struct {
	Items []ChartMapping `json:"items"`
}

func (c *ChartMappingCollection) UnmarshalJSON(b []byte) error {
	switch helpers.JSONType(b) {
	case "ARRAY":
		return json.Unmarshal(b, &c.Items)
	case "OBJECT":
		c.Items = make([]ChartMapping, 1)
		return json.Unmarshal(b, &c.Items[0])
	default:
		return errors.New("payload not recognized as JSON array or object")
	}
}

var listChartsSQL = `SELECT 
							l.slug AS location_slug, 
							v.name, v.slug, v.type_id, 
							p."name" AS provider_name, 
							p.slug AS provider_slug
							FROM chart v
							LEFT JOIN "location" l ON l.id = v.location_id
							JOIN provider p ON p.id = v.provider_id`

func ListCharts(db *pgxpool.Pool) ([]Chart, error) {
	vv := make([]Chart, 0)
	if err := pgxscan.Select(
		context.Background(), db, &vv, listChartsSQL+" ORDER BY p.slug, v.slug",
	); err != nil {
		return make([]Chart, 0), err
	}
	return vv, nil
}

// GetChart returns a single Chart
func GetChart(db *pgxpool.Pool, chartSlug *string) (*Chart, error) {

	var getChartSQL = `
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
								JOIN provider p ON p.id = c.provider_id
								LEFT JOIN chart_variable_mapping cvm ON cvm.chart_id = c.id
								LEFT JOIN timeseries t ON t.id = cvm.timeseries_id
								LEFT JOIN datasource d ON d.id = t.datasource_id 
								LEFT JOIN datatype dt ON dt.id = d.datatype_id
								LEFT JOIN timeseries_providers tp ON tp.id = t.id
								WHERE lower(c.slug) = lower($1)
								GROUP BY l.slug, c.slug, c.name, c.type_id, p.name, p.slug
								LIMIT 1`

	var v Chart
	//if err := pgxscan.Get(context.Background(), db, &v, listChartsSQL+" WHERE lower(v.slug) = lower($1)", chartSlug); err != nil {
	if err := pgxscan.Get(context.Background(), db, &v, getChartSQL, chartSlug); err != nil {
		return nil, err
	}
	return &v, nil
}

func GetChartByLocation(db *pgxpool.Pool, locationSlug *string, chartTypeId *uuid.UUID) (*Chart, error) {

	var getChartSQL = `SELECT
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

	var v Chart
	if err := pgxscan.Get(context.Background(), db, &v, getChartSQL, locationSlug, chartTypeId); err != nil {
		return nil, err
	}
	return &v, nil
}

// CreateChart creates a single Chart
func CreateChart(db *pgxpool.Pool, v *Chart) (*Chart, error) {

	// Insert Into Database Using New Slug
	slug, err := helpers.NextUniqueSlug(db, "chart", "slug", v.Name, "", "")
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

	return GetChart(db, &vSlug)

}

func CreateOrUpdateChartMapping(db *pgxpool.Pool, c ChartMappingCollection, chartSlug *string) ([]ChartMapping, error) {

	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]ChartMapping, 0), err
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
			*chartSlug, v.Variable, v.Key, v.Datatype, v.Provider,
		)
		if err != nil {
			return make([]ChartMapping, 0), err
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

func DeleteChart(db *pgxpool.Pool, chartSlug *string) error {
	if _, err := db.Exec(context.Background(), `DELETE FROM chart WHERE slug=$1`, chartSlug); err != nil {
		return err
	}
	return nil
}
