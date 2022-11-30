package models

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/USACE/water-api/api/helpers"
	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	// TimeseriesGroupMember contains the minimal amount of information
	// necessary to uniquely define a timeseries
	TimeseriesGroupMember struct {
		Provider string `json:"provider"`
		Datatype string `json:"datatype"`
		Key      string `json:"key"`
	}

	TimeseriesGroupMemberCollection struct {
		Items []TimeseriesGroupMember `json:"items"`
	}
)

func (c *TimeseriesGroupMemberCollection) UnmarshalJSON(b []byte) error {
	switch helpers.JSONType(b) {
	case "ARRAY":
		return json.Unmarshal(b, &c.Items)
	case "OBJECT":
		c.Items = make([]TimeseriesGroupMember, 1)
		return json.Unmarshal(b, &c.Items[0])
	default:
		return errors.New("payload not recognized as JSON array or object")
	}
}

func AddTimeseriesGroupMembers(db *pgxpool.Pool, provider *string, timeseriesGroupSlug *string, mc *TimeseriesGroupMemberCollection) (*TimeseriesGroupDetail, error) {

	tx, err := db.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	for _, m := range mc.Items {

		_, err := tx.Exec(
			context.Background(),
			`INSERT INTO timeseries_group_members (timeseries_group_id, timeseries_id) VALUES
			 (
				(SELECT id FROM v_timeseries_group WHERE provider = LOWER($1) AND slug = LOWER($2)),
				(SELECT id FROM v_timeseries WHERE provider=LOWER($3) AND datatype=LOWER($4) AND key=$5)
			 )
			 RETURNING timeseries_group_id`, provider, timeseriesGroupSlug, m.Provider, m.Datatype, m.Key,
		)
		if err != nil {
			tx.Rollback(context.Background())
			return nil, err
		}
	}
	tx.Commit(context.Background())

	return GetTimeseriesGroupDetail(db, &TimeseriesGroupFilter{Slug: timeseriesGroupSlug})

}

func RemoveTimeseriesGroupMembers(db *pgxpool.Pool, provider *string, timeseriesGroupSlug *string, mc *TimeseriesGroupMemberCollection) (*TimeseriesGroupDetail, error) {

	tx, err := db.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	for _, m := range mc.Items {

		_, err := tx.Exec(
			context.Background(),
			`DELETE FROM timeseries_group_members
			 WHERE timeseries_group_id = (SELECT id FROM v_timeseries_group WHERE provider = LOWER($1) AND slug = LOWER($2))
			   AND timeseries_id = (SELECT id FROM v_timeseries WHERE provider=LOWER($3) AND datatype=LOWER($4) AND key=$5)
			 RETURNING timeseries_group_id`, provider, timeseriesGroupSlug, m.Provider, m.Datatype, m.Key,
		)
		if err != nil {
			tx.Rollback(context.Background())
			return nil, err
		}
	}
	tx.Commit(context.Background())

	return GetTimeseriesGroupDetail(db, &TimeseriesGroupFilter{Slug: timeseriesGroupSlug})

}
