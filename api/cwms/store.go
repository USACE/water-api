package cwms

import (
	"github.com/USACE/water-api/api/chartserver"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	Connection  *pgxpool.Pool
	ChartServer *chartserver.ChartServer
}
