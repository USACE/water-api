package main

import (
	"log"
	"net/http"

	"github.com/USACE/water-api/app"
	"github.com/USACE/water-api/cwms"
	"github.com/USACE/water-api/middleware"
	"github.com/USACE/water-api/usgs"

	_ "github.com/jackc/pgx/v4"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
)

func main() {

	// parse configuration from environment variables
	var config app.Config
	if err := envconfig.Process("water", &config); err != nil {
		log.Fatal(err.Error())
	}
	// create store (database pool) from configuration
	st, err := app.NewStore(config)
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
	key := e.Group("")
	key.Use(middleware.KeyAuth(config.ApplicationKey))

	// Health Check
	public.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{"status": "healthy"})
	})

	/////////////////////////////////////////////////////////////////////////////
	// CWMS
	/////////////////////////////////////////////////////////////////////////////
	// CWMS Store
	cs := cwms.Store{Connection: st.Connection}

	// Search
	public.GET("/search/locations", cs.SearchLocations)

	// Locations
	public.GET("/locations", cs.ListLocations)
	public.GET("/locations/:location_id", cs.GetLocation)
	public.GET("/locations/:location_slug/details", cs.GetLocationDetail)
	key.POST("/locations", cs.CreateLocations)
	key.PUT("/locations/:location_id", cs.UpdateLocation)
	key.DELETE("/locations/:location_id", cs.DeleteLocation)

	// Locations (Office Context)
	// public.GET("/offices/:office_symbol/locations")
	// public.GET("/offices/:office_symbol/locations/:location_slug")  // e.g. /offices/lrn/locations/barkley
	key.POST("/offices/:office_symbol/locations", cs.CreateLocationsByOffice)
	key.PUT("/offices/:office_symbol/locations/:location_slug", cs.UpdateLocationByOffice)
	key.DELETE("/offices/:office_symbol/locations/:location_slug", cs.DeleteLocationByOffice)

	// Statistics
	public.GET("/stats/states", cs.ListStatsStates)
	public.GET("/stats/states/:state_id", cs.GetStatsState)
	// public.GET("/stats/watersheds", cs.ListStatsWatersheds)
	// public.GET("/stats/watersheds/:watershed_id", cs.GetStatsWatershed)
	public.GET("/stats/offices", cs.ListStatsOffices)
	public.GET("/stats/offices/:office_id", cs.GetStatsOffice)

	// Sync Locations
	key.POST("/sync/locations", cs.SyncLocations)

	// Location Kinds
	public.GET("/location_kind", cs.ListLocationKind)

	// Offices
	public.GET("/offices", cs.ListOffices)

	// Projects
	public.GET("/projects", cs.ListProjects)

	// States
	public.GET("/states", cs.ListStates)

	// Watersheds
	public.GET("/watersheds", cs.ListWatersheds)
	public.GET("/watersheds/:watershed_id", cs.GetWatershed)
	key.POST("/watersheds", cs.CreateWatershed)

	// Maintenance/Automation
	key.POST("/automation/assign_states_to_locations", cs.AssignStatesToLocations)

	/////////////////////////////////////////////////////////////////////////////
	// USGS
	/////////////////////////////////////////////////////////////////////////////
	// USGS Store
	gs := usgs.Store{Connection: st.Connection}

	// USGS Sites
	public.GET("/usgs/sites", gs.ListSites)
	//public.GET("/usgs/sites/state=:state_abbrev", gs.ListSites)
	public.GET("/usgs/parameters", gs.ListParameters)
	//public.GET("/usgs_sites/enabled_parameters", cs.ListParametersEnabled)

	key.POST("/usgs/sync/sites", gs.SyncSites)
	key.POST("/usgs/site_parameters", gs.CreateSiteParameters)

	// Start server
	log.Fatal(http.ListenAndServe(":80", e))
}
