package models

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

// NwsStages is a NwsStages struct
type NwsStages struct {
	ID             *uuid.UUID `json:"-" db:"id"`
	NwsID          string     `json:"nwsid" db:"nwsid" param:"nwsid"`
	UsgsSiteNumber string     `json:"usgs_site_number" db:"usgs_site_number"`
	Name           string     `json:"name"`
	Stages
}

// Stages is a NwsStages struct
type Stages struct {
	Action        *float64 `json:"action" db:"action"`
	Flood         *float64 `json:"flood" db:"flood"`
	ModerateFlood *float64 `json:"moderate_flood" db:"moderate_flood"`
	MajorFlood    *float64 `json:"major_flood" db:"major_flood"`
}

// WatershedSQL includes common fields selected to build a watershed
const NwsStagesSQL = `SELECT s.id,
                             s.nwsid,
							 s.usgs_site_number,
                             s.name,
							 s.action,
							 s.flood,
							 s.moderate_flood,
							 s.major_flood`

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
		`INSERT INTO nws_stages (name, nwsid, usgs_site_number, action, flood, moderate_flood, major_flood) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, nwsid `,
		&s.Name, &s.NwsID, &s.UsgsSiteNumber, s.Action, s.Flood, s.ModerateFlood, s.MajorFlood,
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
		`UPDATE nws_stages SET name=$1, action=$2, flood=$3, moderate_flood=$4, major_flood=$5, 
		update_date=CURRENT_TIMESTAMP WHERE nwsid=$6 RETURNING nwsid`,
		&s.Name, &s.Action, s.Flood, s.ModerateFlood, s.MajorFlood, s.NwsID); err != nil {
		return nil, err
	}
	return GetNwsStages(db, &nwsID)
}
