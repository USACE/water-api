package models

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

// NwsStages is a NwsStages struct
type NwsStages struct {
	ID                 *uuid.UUID `json:"-" db:"id"`
	NwsID              string     `json:"nwsid" db:"nwsid" param:"nwsid"`
	UsgsSiteNumber     string     `json:"usgs_site_number" db:"usgs_site_number"`
	Name               string     `json:"name"`
	ActionStage        *float64   `json:"action_stage" db:"action_stage"`
	FloodStage         *float64   `json:"flood_stage" db:"flood_stage"`
	ModerateFloodStage *float64   `json:"moderate_flood_stage" db:"moderate_flood_stage"`
	MajorFloodStage    *float64   `json:"major_flood_stage" db:"major_flood_stage"`
}

// WatershedSQL includes common fields selected to build a watershed
const NwsStagesSQL = `SELECT s.id,
                             s.nwsid,
							 s.usgs_site_number,
                             s.name,
							 s.action_stage,
							 s.flood_stage,
							 s.moderate_flood_stage,
							 s.major_flood_stage`

// ListNwsStages returns an array of NWS Location Stages
func ListNwsStages(db *pgxpool.Pool) ([]NwsStages, error) {
	ss := make([]NwsStages, 0)
	if err := pgxscan.Select(context.Background(), db, &ss, NwsStagesSQL+" FROM nws_stages s order by s.name"); err != nil {
		return make([]NwsStages, 0), nil
	}
	return ss, nil
}

// GetNwsStages returns a single NWS Stages Record
func GetNwsStages(db *pgxpool.Pool, NwsOrUsgsID *string) (*NwsStages, error) {
	var s NwsStages
	if err := pgxscan.Get(
		context.Background(), db, &s, NwsStagesSQL+` FROM nws_stages s WHERE s.nwsid = $1 OR s.usgs_site_number = $2`, NwsOrUsgsID, NwsOrUsgsID,
	); err != nil {
		return nil, err
	}
	return &s, nil
}

// CreateNwsStage creates a new NWS Stages Record
func CreateNwsStage(db *pgxpool.Pool, s *NwsStages) (*NwsStages, error) {
	//slug, err := helpers.NextUniqueSlug(db, "watershed", "slug", w.Name, "", "")
	// if err != nil {
	// 	return nil, err
	// }
	var sNew NwsStages
	if err := pgxscan.Get(context.Background(), db, &sNew,
		`INSERT INTO nws_stages (name, nwsid, usgs_site_number, action_stage, flood_stage, moderate_flood_stage, major_flood_stage) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, nwsid `,
		&s.Name, &s.NwsID, &s.UsgsSiteNumber, s.ActionStage, s.FloodStage, s.ModerateFloodStage, s.MajorFloodStage,
	); err != nil {
		return nil, err
	}
	return GetNwsStages(db, &sNew.NwsID)
	//return &wNew, nil
}

// UpdateNwsStages updates a single NWS Stages Record
func UpdateNwsStages(db *pgxpool.Pool, s *NwsStages) (*NwsStages, error) {
	var nwsID string
	if err := pgxscan.Get(context.Background(), db, &nwsID,
		`UPDATE nws_stages SET name=$1, action_stage=$2, flood_stage=$3, moderate_flood_stage=$4, major_flood_stage=$5, 
		update_date=CURRENT_TIMESTAMP WHERE nwsid=$6 RETURNING nwsid`,
		&s.Name, &s.ActionStage, s.FloodStage, s.ModerateFloodStage, s.MajorFloodStage, s.NwsID); err != nil {
		return nil, err
	}
	return GetNwsStages(db, &nwsID)
}
