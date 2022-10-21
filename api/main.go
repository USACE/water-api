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
	"github.com/USACE/water-api/api/providers"
	"github.com/USACE/water-api/api/ratings"
	"github.com/USACE/water-api/api/usgs"
	"github.com/USACE/water-api/api/visualizations"
	"github.com/USACE/water-api/api/watersheds"

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

	// d3-chart-server integration
	chartserver, err := chartserver.NewChartServer(chartserver.Config{URLString: config.ChartServerURL})
	if err != nil {
		log.Fatal(err.Error())
	}

	// All Routes
	e := echo.New()                         // All Routes
	e.Use(middleware.CORS, middleware.GZIP) // All Routes Middleware

	// Public Routes
	public := e.Group("") // Public Routes

	// Public Health Check Route
	public.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{"status": "healthy"})
	})

	// pg-featureserv routes
	features := public.Group("/features")
	features.Use(middleware.PgFeatureservProxy(config.PgFeatureservUrl))

	// Mount Routes
	cwms.Mount(st.Connection, e, &config, chartserver) // CWMS
	nws.Mount(st.Connection, e, &config)               // National Weather Service
	providers.Mount(st.Connection, e, &config)         // Providers
	usgs.Mount(st.Connection, e, &config)              // USGS
	visualizations.Mount(st.Connection, e, &config)    // Visualizations
	watersheds.Mount(st.Connection, e, &config)        // Watersheds
	ratings.Mount(st.Connection, e, &config)           // Ratings

	// Start Server
	s := &http2.Server{
		MaxConcurrentStreams: 250,     // http2 default 250
		MaxReadFrameSize:     1048576, // http2 default 1048576
		IdleTimeout:          10 * time.Second,
	}
	if err := e.StartH2CServer(":80", s); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
