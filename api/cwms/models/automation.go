package models

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func AssignStatesToLocations(db *pgxpool.Pool) error {

	if _, err := db.Exec(
		context.Background(),
		`WITH location_states AS (
	    	SELECT l.id  AS location_id,
	    		   s.gid AS state_id
	    	FROM location l
	    	JOIN tiger_data.state_all s ON ST_Contains(s.the_geom, ST_Transform(l.geometry, 4269))
	    )
	    UPDATE location l
	    SET state_id = (SELECT state_id FROM location_states WHERE location_states.location_id = l.id)`,
	); err != nil {
		return err
	}

	return nil
}
