package locations

import (
	"github.com/USACE/water-api/api/app"
	"github.com/USACE/water-api/api/chartserver"
	"github.com/USACE/water-api/api/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

type Store struct {
	Connection  *pgxpool.Pool
	ChartServer *chartserver.ChartServer
}

func Mount(conn *pgxpool.Pool, e *echo.Echo, config *app.Config, chartserver *chartserver.ChartServer) {

	s := Store{
		Connection:  conn,
		ChartServer: chartserver,
	}

	// Public Routes
	public := e.Group("")

	// Key Only Group
	key := e.Group("")
	key.Use(middleware.KeyAuth(config.ApplicationKey))

	// LOCATIONS; :location corresponds to unique location slug
	public.GET("/v2/locations", s.ListLocations)                           // LIST
	public.GET("/v2/locations/:location", s.GetLocation)                   // GET
	public.GET("/v2/locations/:location/details", s.GetLocationDetail)     // GET DETAILS
	public.GET("/v2/locations/:location/profile-chart", s.GetProfileChart) // GET PROFILE CHART
	key.POST("/v2/locations", s.CreateLocations)                           // CREATE
	key.PUT("/v2/locations/:location", s.UpdateLocation)                   // UPDATE
	key.DELETE("/v2/locations/:location", s.DeleteLocation)                // DELETE

	// Sync Locations
	key.POST("/v2/sync/locations", s.SyncLocations)
}