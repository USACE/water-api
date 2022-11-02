package cwms

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

	// LOCATION KINDS
	public.GET("/location_kind", s.ListLocationKind)

	// LOCATIONS
	public.GET("/search/locations", s.SearchLocations)                       // SEARCH
	public.GET("/locations", s.ListLocations)                                // LIST
	public.GET("/locations/:location_id", s.GetLocation)                     // GET
	public.GET("/locations/:location_slug/details", s.GetLocationDetail)     // GET DETAILS
	public.GET("/locations/:location_slug/profile-chart", s.GetProfileChart) // GET PROFILE CHART
	key.POST("/locations", s.CreateLocations)                                // CREATE
	key.PUT("/locations/:location_id", s.UpdateLocation)                     // UPDATE
	key.DELETE("/locations/:location_id", s.DeleteLocation)                  // DELETE

	// LOCATIONS (Office Context)
	// public.GET("/offices/:office_symbol/locations")
	// public.GET("/offices/:office_symbol/locations/:location_slug")  // e.g. /offices/lrn/locations/barkley
	key.POST("/offices/:office_symbol/locations", s.CreateLocationsByOffice)
	key.PUT("/offices/:office_symbol/locations/:location_slug", s.UpdateLocationByOffice)
	key.DELETE("/offices/:office_symbol/locations/:location_slug", s.DeleteLocationByOffice)

	// LOCATION LEVEL KINDS
	public.GET("/levels/kind", s.ListLevelKind)         // LIST
	key.POST("/levels/kind/:name", s.CreateLevelKind)   // CREATE
	key.DELETE("/levels/kind/:slug", s.DeleteLevelKind) // DELETE

	// LOCATION LEVELS
	public.GET("/levels/:location_slug/:level_kind", s.ListLevelValues)
	key.POST("/levels/:location_id", s.CreateLocationLevels)
	key.PUT("/levels/:location_id", s.UpdateLocationLevels)

	// Statistics
	public.GET("/stats/states", s.ListStatsStates)
	public.GET("/stats/states/:state_id", s.GetStatsState)
	public.GET("/stats/offices", s.ListStatsOffices)
	public.GET("/stats/offices/:office_id", s.GetStatsOffice)

	// Sync Locations
	key.POST("/sync/locations", s.SyncLocations)

	// Offices
	public.GET("/offices", s.ListOffices)

	// Projects
	public.GET("/projects", s.ListProjects)

	// States
	public.GET("/states", s.ListStates)

	// // Timeseries
	// public.GET("/timeseries", s.ListTimeseries)
	// key.POST("/timeseries/measurements", s.CreateOrUpdateTimeseriesMeasurements)
	// key.POST("/timeseries", s.CreateTimeseries) // (airflow - array of objects in payload)
	// public.POST "/:provider_slug/timeseries"
	// "/levels/latest/config/:owner"

	// Maintenance/Automation
	key.POST("/automation/assign_states_to_locations", s.AssignStatesToLocations)
}
