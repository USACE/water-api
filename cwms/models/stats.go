package models

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ListStatsStatesBaseSQL = `WITH
	state_location_counts AS (
		SELECT l.state_id, COUNT(*)
		FROM a2w_cwms.location l
		WHERE l.state_id IS NOT null
		GROUP BY l.state_id
	),
		state_project_counts AS (
		SELECT l.state_id, COUNT(*)
		FROM a2w_cwms.location l
		WHERE l.state_id IS NOT null AND l.kind_id = (SELECT id FROM a2w_cwms.location_kind lk WHERE lk.name = 'PROJECT')
		GROUP BY l.state_id
	)
	SELECT l.state_id,
		l.count AS locations,
		p.count AS projects
	FROM state_location_counts l
	JOIN state_project_counts p ON p.state_id = l.state_id`

var ListStatsOfficeBaseSQL = `WITH
	office_locations AS (
		SELECT l.office_id, count(*) AS locations
		FROM a2w_cwms.location l
		GROUP BY l.office_id
	)
	SELECT office_locations.office_id, office_locations.locations
	FROM office_locations`

type StatsState struct {
	StateID   int8  `json:"state_id"`
	Locations int16 `json:"locations"`
	Projects  int16 `json:"projects"`
}

type StatsOffice struct {
	OfficeID  uuid.UUID `json:"office_id"`
	Locations int16     `json:"locations"`
}

func ListStatsStates(db *pgxpool.Pool) ([]StatsState, error) {
	ss := make([]StatsState, 0)
	if err := pgxscan.Select(
		context.Background(),
		db,
		&ss,
		ListStatsStatesBaseSQL,
	); err != nil {
		return make([]StatsState, 0), err
	}
	return ss, nil
}

func GetStatsState(db *pgxpool.Pool, id *string) (*StatsState, error) {
	var ss StatsState
	if err := pgxscan.Get(
		context.Background(),
		db,
		&ss,
		ListStatsStatesBaseSQL+" WHERE l.state_id = $1",
		id,
	); err != nil {
		return nil, err
	}
	return &ss, nil
}

func ListStatsOffices(db *pgxpool.Pool) ([]StatsOffice, error) {
	so := make([]StatsOffice, 0)
	if err := pgxscan.Select(
		context.Background(),
		db,
		&so,
		ListStatsOfficeBaseSQL,
	); err != nil {
		return make([]StatsOffice, 0), err
	}
	return so, nil
}

func GetStatsOffice(db *pgxpool.Pool, officeID *uuid.UUID) (*StatsOffice, error) {
	var so StatsOffice
	if err := pgxscan.Get(
		context.Background(),
		db,
		&so,
		ListStatsOfficeBaseSQL+" WHERE office_locations.office_id = $1",
		officeID,
	); err != nil {
		return nil, err
	}
	return &so, nil
}
