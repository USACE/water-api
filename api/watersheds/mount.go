package watersheds

import (
	"log"

	"github.com/USACE/water-api/api/app"
	"github.com/USACE/water-api/api/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

type Store struct {
	Connection *pgxpool.Pool
}

func Mount(conn *pgxpool.Pool, e *echo.Echo, config *app.Config) {

	s := Store{Connection: conn} // database connection

	// PUBLIC ROUTES
	public := e.Group("")

	// APPKEY ROUTES (AUTHENTICATED BY APPLICATION KEY ONLY)
	key := e.Group("")
	key.Use(middleware.KeyAuth(config.ApplicationKey))

	// PRIVATE ROUTES (JWT)
	private := e.Group("")
	if config.AuthMocked {
		log.Println("Auth is Disabled...") // @todo; re-add middleware.JWTMock
	} else {
		private.Use(middleware.JWT, middleware.AttachUserInfo)
	}

	// WATERSHED CRUD OPERATIONS
	public.GET("/watersheds", s.ListWatersheds)                               // LIST
	public.GET("/watersheds/:watershed_slug", s.GetWatershed)                 // GET
	private.POST("/watersheds", s.CreateWatershed)                            // CREATE
	private.PUT("/watersheds/:watershed_slug", s.UpdateWatershed)             // UPDATE
	private.DELETE("/watersheds/:watershed_slug", s.DeleteWatershed)          // SOFT DELETE
	private.POST("/watersheds/:watershed_slug/undelete", s.UndeleteWatershed) // UNDELETE

	// WATERSHED GEOMETRY OPERATIONS
	public.GET("/watersheds/:watershed_slug/geometry", s.GetWatershedGeometry)      // GET GEOJSON GEOMETRY
	key.PUT("/watersheds/:watershed_id/update_geometry", s.UpdateWatershedGeometry) // EDIT GEOJSON GEOMETRY

	// WATERSHED SHAPEFILE UPLOAD
	private.POST("/watersheds/:watershed_id/shapefile_uploads", s.UploadWatersheds, middleware.IsAdmin) // UPLOAD SHAPEFILE

	// WATERSHED DATA EXTRACT OPERATIONS; @todo; This may belong in another module
	// public.GET("watersheds/:watershed_slug/extract", s.WatershedExtract) // Extract Timeseries values for all locations associated with a watershed

	// ASSOCIATE USGS SITES/PARAMETERS WITH A WATERSHED
	public.GET("/watersheds/usgs_sites", s.ListWatershedSiteParameters)                                                       // List watershed USGS Site Params mapped for data retrieval; Used by Airflow.
	private.POST("/watersheds/:watershed_slug/site/:site_number/parameter/:parameter_code", s.CreateWatershedSiteParameter)   // Associate site parameter with watershed
	private.DELETE("/watersheds/:watershed_slug/site/:site_number/parameter/:parameter_code", s.DeleteWatershedSiteParameter) // Un-Associate site parameter with watershed

	// STATS ENDPOINTS @todo
	// public.GET("/stats/watersheds", cs.ListStatsWatersheds)
	// public.GET("/stats/watersheds/:watershed_id", cs.GetStatsWatershed)
}
