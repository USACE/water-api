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

type Timeseries struct {
	ID             *uuid.UUID `json:"id"`
	ProviderID     *string    `json:"provider_id"`
	DatasourceType string     `json:"datasource_type" db:"datasource_type"`
	Key            string     `json:"key"`
	LatestTime     time.Time  `json:"latest_time" db:"latest_time"`
	LatestValue    float64    `json:"latest_value" db:"latest_value"`
}

type TimeseriesFilter struct {
	KindID         *uuid.UUID `json:"kind_id" query:"kind_id"`
	DatasourceType string     `json:"datasource_type" query:"datasource_type"`
	ProviderID     *uuid.UUID `json:"provider_id" query:"provider_id"`
	Q              *string    `query:"q"`
}

type TimeseriesCollection struct {
	Items []Timeseries `json:"items"`
}

func (c *TimeseriesCollection) UnmarshalJSON(b []byte) error {
	switch helpers.JSONType(b) {
	case "ARRAY":
		return json.Unmarshal(b, &c.Items)
	case "OBJECT":
		c.Items = make([]Timeseries, 1)
		return json.Unmarshal(b, &c.Items[0])
	default:
		return errors.New("payload not recognized as JSON array or object")
	}
}

func SyncTimeseries(db *pgxpool.Pool, c TimeseriesCollection) ([]Timeseries, error) {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]Timeseries, 0), err
	}
	defer tx.Rollback(context.Background())

	newIDs := make([]uuid.UUID, 0)

	for _, t := range c.Items {
		rows, err := tx.Query(
			context.Background(),
			`INSERT INTO timeseries (datasource_id, datasource_key, latest_time, latest_value)
			VALUES((SELECT d.id FROM a2w_cwms.datasource d 
				JOIN a2w_cwms.datasource_type dt ON dt.id = d.datasource_type_id 
				WHERE dt.slug = $1), $2, $3, $4)
			ON CONFLICT ON CONSTRAINT timeseries_unique_datasource
			DO UPDATE SET
			latest_time = $3,
			latest_value = $4
			WHERE timeseries.datasource_key = $2
			--AND $3 >= timeseries.latest_time
			RETURNING id`,
			t.DatasourceType, t.Key, t.LatestTime, t.LatestValue,
		)
		if err != nil {
			return make([]Timeseries, 0), err
		}
		var id uuid.UUID
		if err := pgxscan.ScanOne(&id, rows); err != nil {
			tx.Rollback(context.Background())
			return c.Items, err
		} else {
			newIDs = append(newIDs, id)
		}
	}
	tx.Commit(context.Background())
	return nil, nil
}
