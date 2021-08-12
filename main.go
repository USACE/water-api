package main

import (
	"log"
	"net/http"

	"github.com/USACE/water-api/cwms"
	"github.com/USACE/water-api/middleware"

	_ "github.com/jackc/pgx/v4"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
)

func main() {

	// parse configuration from environment variables
	var config cwms.Config
	if err := envconfig.Process("water", &config); err != nil {
		log.Fatal(err.Error())
	}
	// create cwms store (database pool) from configuration
	cs, err := cwms.NewStore(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	e := echo.New()                         // All Routes
	e.Use(middleware.CORS, middleware.GZIP) // All Routes Middleware

	// Public Routes
	public := e.Group("")

	// Private Routes w/ Access Control
	private := e.Group("")
	if config.AuthMocked {
		private.Use(middleware.JWTMock)
	} else {
		private.Use(middleware.JWT)
	}

	// App Routes (Intended to be used by application only)
	app := e.Group("")
	app.Use(middleware.KeyAuth(config.ApplicationKey))

	// Health Check
	public.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{"status": "healthy"})
	})

	// Search
	public.GET("/search/locations", cs.SearchLocations)

	// Locations
	public.GET("/locations", cs.ListLocations)
	app.POST("/locations", cs.CreateLocations)
	public.GET("/locations/:location_id", cs.GetLocation)
	app.PUT("/locations/:location_id", cs.UpdateLocation)
	app.DELETE("/locations/:location_id", cs.DeleteLocation)

	// Locations (Office Context)
	// public.GET("/offices/:office_symbol/locations")
	// public.GET("/offices/:office_symbol/locations/:location_slug")  // e.g. /offices/lrn/locations/barkley
	private.POST("/offices/:office_symbol/locations", cs.CreateLocationsByOffice)
	private.PUT("/offices/:office_symbol/locations/:location_slug", cs.UpdateLocationByOffice)
	private.DELETE("/offices/:office_symbol/locations/:location_slug", cs.DeleteLocationByOffice)

	// Statistics
	public.GET("/stats/states", cs.ListStatsStates)
	public.GET("/stats/states/:state_id", cs.GetStatsState)
	// public.GET("/stats/watersheds", cs.ListStatsWatersheds)
	// public.GET("/stats/watersheds/:watershed_id", cs.GetStatsWatershed)
	public.GET("/stats/offices", cs.ListStatsOffices)
	public.GET("/stats/offices/:office_id", cs.GetStatsOffice)

	// Sync Postgres
	app.POST("/sync/locations", cs.SyncLocations)

	// Location Kinds
	public.GET("/location_kind", cs.ListLocationKind)

	// Projects
	public.GET("/projects", cs.ListProjects)

	// States
	public.GET("/states", cs.ListStates)

	// USGS Sites
	public.GET("/usgs_sites", cs.ListSites)
	public.GET("/usgs_sites/state/:state_abbrev", cs.ListSites)
	restricted.POST("/usgs_sites", cs.SyncSites)

	// Maintenance/Automation
	app.POST("/automation/assign_states_to_locations", cs.AssignStatesToLocations)

	// Start server
	log.Fatal(http.ListenAndServe(":80", e))
}
