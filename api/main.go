package main

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/net/http2"

	"github.com/USACE/water-api/api/app"
	"github.com/USACE/water-api/api/chartserver"
	"github.com/USACE/water-api/api/cwms"
	"github.com/USACE/water-api/api/middleware"
	"github.com/USACE/water-api/api/nws"
	"github.com/USACE/water-api/api/usgs"
	"github.com/USACE/water-api/api/water"

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
		// @todo. re-add JWTMock
		// private.Use(middleware.JWTMock)
		log.Println("Auth is Disabled...")
	} else {
		private.Use(middleware.JWT, middleware.AttachUserInfo)
	}

	// Routes to Serve Feature Data
	features := public.Group("/features")
	features.Use(middleware.PgFeatureservProxy(config.PgFeatureservUrl))

	// App Routes (Intended to be used by application only)
	key := e.Group("")
	key.Use(middleware.KeyAuth(config.ApplicationKey))

	// Health Check
	public.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{"status": "healthy"})
	})

	/////////////////////////////////////////////////////////////////////////////
	// CHARTSERVER INTEGRATION
	/////////////////////////////////////////////////////////////////////////////
	chartserver, err := chartserver.NewChartServer(chartserver.Config{URLString: config.ChartServerURL})
	if err != nil {
		log.Fatal(err.Error())
	}

	/////////////////////////////////////////////////////////////////////////////
	// CWMS
	/////////////////////////////////////////////////////////////////////////////
	// CWMS Store
	cs := cwms.Store{
		Connection:  st.Connection,
		ChartServer: chartserver,
	}

	// Search
	public.GET("/search/locations", cs.SearchLocations)

	// Locations
	public.GET("/locations", cs.ListLocations)
	public.GET("/locations/:location_id", cs.GetLocation)
	public.GET("/locations/:location_slug/details", cs.GetLocationDetail)
	key.POST("/locations", cs.CreateLocations)
	key.PUT("/locations/:location_id", cs.UpdateLocation)
	key.DELETE("/locations/:location_id", cs.DeleteLocation)

	// Get Location Profile Chart
	public.GET("/locations/:location_id/profile-chart", cs.GetProfileChart)

	// Locations (Office Context)
	// public.GET("/offices/:office_symbol/locations")
	// public.GET("/offices/:office_symbol/locations/:location_slug")  // e.g. /offices/lrn/locations/barkley
	key.POST("/offices/:office_symbol/locations", cs.CreateLocationsByOffice)
	key.PUT("/offices/:office_symbol/locations/:location_slug", cs.UpdateLocationByOffice)
	key.DELETE("/offices/:office_symbol/locations/:location_slug", cs.DeleteLocationByOffice)

	// Location Levels
	public.GET("/levels/kind", cs.ListLevelKind)
	key.POST("/levels/kind/:name", cs.CreateLevelKind)
	key.DELETE("/levels/kind/:slug", cs.DeleteLevelKind)
	public.GET("/levels/:location_slug/:level_kind", cs.ListLevelValues)
	key.POST("/levels/:location_id", cs.CreateLocationLevels)
	key.PUT("/levels/:location_id", cs.UpdateLocationLevels)

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
	public.GET("/watersheds/:watershed_slug", cs.GetWatershed)
	private.POST("/watersheds", cs.CreateWatershed)
	private.PUT("/watersheds/:watershed_slug", cs.UpdateWatershed)
	public.GET("/watersheds/:watershed_slug/geometry", cs.GetWatershedGeometry)
	key.PUT("/watersheds/:watershed_id/update_geometry", cs.UpdateWatershedGeometry)
	private.DELETE("/watersheds/:watershed_slug", cs.DeleteWatershed)
	private.POST("/watersheds/:watershed_slug/undelete", cs.UndeleteWatershed)
	private.POST("/watersheds/:watershed_id/shapefile_uploads", cs.UploadWatersheds, middleware.IsAdmin)

	// Extract timeseries values using locations grouped with a watershed defined by its slug
	public.GET("watersheds/:watershed_slug/extract", cs.WatershedExtract)

	// Maintenance/Automation
	key.POST("/automation/assign_states_to_locations", cs.AssignStatesToLocations)

	/////////////////////////////////////////////////////////////////////////////
	// USGS
	/////////////////////////////////////////////////////////////////////////////
	// USGS Store
	gs := usgs.Store{Connection: st.Connection}

	// Search
	public.GET("/search/usgs_sites", gs.SearchSites)

	// USGS Sites
	public.GET("/usgs/sites", gs.ListSites) // Will accept ?state=xx
	public.GET("/usgs/sites/:site_number", gs.GetSite)
	public.GET("/usgs/parameters", gs.ListParameters)
	//public.GET("/usgs_sites/enabled_parameters", cs.ListParametersEnabled)
	key.POST("/usgs/sync/sites", gs.SyncSites)
	key.POST("/usgs/site_parameters", gs.CreateSiteParameters)

	// USGS Time Series
	key.POST("/usgs/sites/:site_number/measurements", gs.CreateOrUpdateUSGSMeasurements)
	public.GET("/usgs/sites/:site_number/measurements", gs.ListUSGSMeasurements)

	/////////////////////////////////////////////////////////////////////////////
	// NWS
	/////////////////////////////////////////////////////////////////////////////
	// NWS Store
	ns := nws.Store{Connection: st.Connection}
	public.GET("/nws/stages", ns.ListNwsStages)
	public.GET("/nws/stages/:nwsid", ns.GetNwsStages)
	key.POST("/nws/stages", ns.CreateNwsStage)
	key.PUT("/nws/stages/:nwsid", ns.UpdateNwsStages)

	/////////////////////////////////////////////////////////////////////////////
	// WATER
	/////////////////////////////////////////////////////////////////////////////
	// WATER Store
	ws := water.Store{Connection: st.Connection}

	// Associate USGS sites/parameters with Watershed
	private.POST("/watersheds/:watershed_slug/site/:site_number/parameter/:parameter_code", ws.CreateWatershedSiteParameter)
	private.DELETE("/watersheds/:watershed_slug/site/:site_number/parameter/:parameter_code", ws.DeleteWatershedSiteParameter)
	// Watershed USGS Site Params enabled for data retrieval.  Primarily used by Airflow.
	public.GET("/watersheds/usgs_sites", ws.ListWatershedSiteParameters)

	// GET "/timeseries?provider=LRH&datasource_type=cwms-timeseries"
	key.POST("/timeseries", cs.SyncTimeseries) // (airflow - array of objects in payload)
	// public.POST "/:provider_slug/timeseries"
	// "/levels/latest/config/:owner"

	// Server
	s := &http2.Server{
		MaxConcurrentStreams: 250,     // http2 default 250
		MaxReadFrameSize:     1048576, // http2 default 1048576
		IdleTimeout:          10 * time.Second,
	}
	if err := e.StartH2CServer(":80", s); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
