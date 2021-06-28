package models

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ListStatsStatesBaseSQL = `WITH state_location_counts AS (
						          SELECT state_id, COUNT(*)
	                         	  FROM location
	                         	  WHERE state_id IS NOT null
	                         	  GROUP BY state_id
	                          ),
							  state_project_counts AS (
								  SELECT state_id, COUNT(*)
								  FROM location
								  WHERE state_id IS NOT null
								  	AND kind_id = (SELECT id FROM location_kind WHERE UPPER(name) = 'PROJECT')
								  GROUP BY state_id
	                         )
	                         SELECT l.state_id,
							        l.count AS locations,
									p.count AS projects
	                         FROM state_location_counts l
	                         JOIN state_project_counts p ON p.state_id = l.state_id`

var ListStatsOfficeBaseSQL = `WITH office_location_counts AS (
	                              SELECT office_id, COUNT(*)
								  FROM location
								  WHERE office_id IS NOT null
								  GROUP BY office_id
							 ),
							 office_project_counts AS (
								SELECT office_id, COUNT(*)
								FROM location
								WHERE state_id IS NOT null
								  AND kind_id = (SELECT id FROM location_kind WHERE UPPER(name) = 'PROJECT')
								GROUP BY office_id
							 )
							 SELECT l.office_id,
							        l.count AS locations,
									p.count AS projects
							 FROM office_location_counts l
							 JOIN office_project_counts p ON p.office_id = l.office_id`

type StatsState struct {
	StateID   int `json:"state_id"`
	Locations int `json:"locations"`
	Projects  int `json:"projects"`
}

type StatsOffice struct {
	OfficeID  uuid.UUID `json:"office_id"`
	Locations int       `json:"locations"`
	Projects  int       `json:"projects"`
}

func ListStatsStates(db *pgxpool.Pool) ([]StatsState, error) {
	ss := make([]StatsState, 0)
	if err := pgxscan.Select(context.Background(), db, &ss, ListStatsStatesBaseSQL); err != nil {
		return make([]StatsState, 0), err
	}
	return ss, nil
}

func GetStatsState(db *pgxpool.Pool, id *string) (*StatsState, error) {
	var ss StatsState
	if err := pgxscan.Get(
		context.Background(), db, &ss,
		ListStatsStatesBaseSQL+" WHERE l.state_id = $1", id,
	); err != nil {
		return nil, err
	}
	return &ss, nil
}

func ListStatsOffices(db *pgxpool.Pool) ([]StatsOffice, error) {
	so := make([]StatsOffice, 0)
	if err := pgxscan.Select(context.Background(), db, &so, ListStatsOfficeBaseSQL); err != nil {
		return make([]StatsOffice, 0), err
	}
	return so, nil
}

func GetStatsOffice(db *pgxpool.Pool, officeID *uuid.UUID) (*StatsOffice, error) {
	var so StatsOffice
	if err := pgxscan.Get(
		context.Background(), db, &so,
		ListStatsOfficeBaseSQL+" WHERE l.office_id = $1", officeID,
	); err != nil {
		return nil, err
	}
	return &so, nil
}
