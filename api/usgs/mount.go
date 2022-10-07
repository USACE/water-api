package usgs

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

	// Database Connection
	s := Store{Connection: conn}

	// Public Routes
	public := e.Group("")

	// Key Only Group
	key := e.Group("")
	key.Use(middleware.KeyAuth(config.ApplicationKey))

	// Search
	public.GET("/search/usgs_sites", s.SearchSites)

	// USGS Sites
	public.GET("/usgs/sites", s.ListSites) // Will accept ?state=xx
	public.GET("/usgs/sites/:site_number", s.GetSite)
	public.GET("/usgs/parameters", s.ListParameters)

	//public.GET("/usgs_sites/enabled_parameters", cs.ListParametersEnabled)
	key.POST("/usgs/sync/sites", s.SyncSites)
	key.POST("/usgs/site_parameters", s.CreateSiteParameters)

	// USGS Time Series
	key.POST("/usgs/sites/:site_number/measurements", s.CreateOrUpdateUSGSMeasurements)
	public.GET("/usgs/sites/:site_number/measurements", s.ListUSGSMeasurements)
}
