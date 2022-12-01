package charts

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/USACE/water-api/api/helpers"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	ChartMapping struct {
		Variable    string         `json:"variable,omitempty" db:"variable"`
		Key         string         `json:"key,omitempty"`
		Datatype    string         `json:"datatype" query:"datatype"`
		Provider    string         `json:"provider"`
		LatestValue *[]interface{} `json:"latest_value" db:"latest_value"`
	}

	ChartMappingCollection struct {
		Items []ChartMapping `json:"items"`
	}
)

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

func CreateOrUpdateChartMapping(db *pgxpool.Pool, chartProvider *string, chartSlug *string, mc *ChartMappingCollection) error {

	tx, err := db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	updatedIDs := make([]uuid.UUID, 0)

	for _, m := range mc.Items {

		rows, err := tx.Query(
			context.Background(),
			`INSERT into chart_variable_mapping(chart_id, variable, timeseries_id)
			VALUES
			(
				(SELECT id FROM v_chart WHERE provider = LOWER($1) AND slug = LOWER($2)),
				$3,
				(SELECT id FROM v_timeseries WHERE provider = LOWER($4) AND datatype = LOWER($5) AND LOWER(key) = LOWER($6))
			)
			ON CONFLICT ON CONSTRAINT chart_unique_variable DO UPDATE SET timeseries_id = EXCLUDED.timeseries_id
			RETURNING chart_id`, chartProvider, chartSlug, m.Variable, m.Provider, m.Datatype, m.Key,
		)
		if err != nil {
			tx.Rollback(context.Background())
			return err
		}
		var id uuid.UUID
		if err := pgxscan.ScanOne(&id, rows); err != nil {
			tx.Rollback(context.Background())
			return err
		}

		updatedIDs = append(updatedIDs, id)
	}

	tx.Commit(context.Background())

	if len(updatedIDs) == 0 {
		return errors.New("no records updated")
	}

	return nil

}

func DeleteChartMapping(db *pgxpool.Pool, chartProvider *string, chartSlug *string, mc *ChartMappingCollection) error {

	tx, err := db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	updatedIDs := make([]uuid.UUID, 0)

	for _, m := range mc.Items {

		rows, err := tx.Query(
			context.Background(),
			`DELETE FROM chart_variable_mapping
			 WHERE chart_id = (SELECT id FROM v_chart WHERE provider = LOWER($1) AND slug = LOWER($2))
			   AND variable = $3
			 RETURNING chart_id`, chartProvider, chartSlug, m.Variable,
		)
		if err != nil {
			tx.Rollback(context.Background())
			return err
		}
		var id uuid.UUID
		if err := pgxscan.ScanOne(&id, rows); err != nil {
			tx.Rollback(context.Background())
			return err
		}

		updatedIDs = append(updatedIDs, id)
	}

	tx.Commit(context.Background())

	if len(updatedIDs) == 0 {
		return errors.New("no records updated")
	}

	return nil

}
