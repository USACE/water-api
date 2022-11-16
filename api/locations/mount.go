package locations

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

	s := Store{Connection: conn}

	// Public Routes
	public := e.Group("")

	// Key Only Group
	key := e.Group("")
	key.Use(middleware.KeyAuth(config.ApplicationKey))

	// LOCATIONS; :location corresponds to unique location slug
	public.GET("/locations", s.ListLocations)                                   // LIST
	public.GET("/locations/:location", s.GetLocation)                           // GET ONE
	key.POST("/providers/:provider/locations", s.CreateLocations)               // CREATE
	key.PUT("/providers/:provider/locations", s.UpdateLocations)                // UPDATE
	key.DELETE("/providers/:provider/locations", s.DeleteLocations)             // DELETE MANY
	key.DELETE("/providers/:provider/locations/:location", s.DeleteOneLocation) // DELETE ONE (USING SLUG)
}
