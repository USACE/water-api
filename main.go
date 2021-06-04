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

	e := echo.New()           // All Routes
	e.Use(middleware.GZIP)    // All Routes Middleware
	public := e.Group("")     // Public Routes
	restricted := e.Group("") // Restricted Routes
	restricted.Use(middleware.KeyAuth(config.ApplicationKey))

	// Locations
	public.GET("/locations", cs.ListLocations)
	restricted.POST("/locations", cs.CreateLocations)
	public.GET("/locations/:location_id", cs.GetLocation)
	restricted.PUT("/locations/:location_id", cs.UpdateLocation)
	restricted.DELETE("/locations/:location_id", cs.DeleteLocation)

	// Projects
	public.GET("/projects", cs.ListProjects)

	// Start server
	log.Fatal(http.ListenAndServe(":80", e))
}
