package nws

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	Connection *pgxpool.Pool
}
