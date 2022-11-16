package timeseries

import (
	"github.com/USACE/water-api/api/app"
	"github.com/USACE/water-api/api/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

type Store struct {
	Connection *pgxpool.Pool
}

func Mount(conn *pgxpool.Pool, e *echo.Echo, config *app.Config) {

	s := Store{
		Connection: conn,
	}

	// Public Routes
	public := e.Group("")

	// Key Only Group
	key := e.Group("")
	key.Use(middleware.KeyAuth(config.ApplicationKey))

	// Timeseries
	public.GET("/timeseries", s.ListTimeseries)
	key.POST("/providers/:provider/timeseries", s.CreateTimeseries) // (airflow - array of objects in payload)
	key.POST("/providers/:provider/timeseries/measurements", s.CreateOrUpdateTimeseriesMeasurements)

}
