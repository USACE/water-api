package locations

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

func UpdateLocation(db *pgxpool.Pool, l *Location) (*Location, error) {
	// var s string
	// if err := pgxscan.Get(
	// 	context.Background(), db, &id,
	// 	"UPDATE location SET update_date=CURRENT_TIMESTAMP, office_id=$2, name=$3, public_name=$4, geometry=$5, kind_id=$6 WHERE id = $1 RETURNING slug",
	// 	l.ID, l.OfficeID, l.Name, l.PublicName, l.Geometry.EWKT(6), l.KindID,
	// ); err != nil {
	// 	return nil, err
	// }
	// return GetLocation(db, &Filter{Slug: &s})
	return nil, nil
}
