package main

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/net/http2"

	"github.com/USACE/water-api/api/app"
	"github.com/USACE/water-api/api/charts"
	"github.com/USACE/water-api/api/locations"
	"github.com/USACE/water-api/api/middleware"
	"github.com/USACE/water-api/api/providers"
	"github.com/USACE/water-api/api/timeseries"

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
	charts.Mount(st.Connection, e, &config)     // Charts
	locations.Mount(st.Connection, e, &config)  // Locations
	providers.Mount(st.Connection, e, &config)  // Providers
	timeseries.Mount(st.Connection, e, &config) // Timeseries

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
