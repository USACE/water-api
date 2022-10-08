package nws

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

	public.GET("/nws/stages", s.ListNwsStages)       // LIST
	public.GET("/nws/stages/:nwsid", s.GetNwsStages) // GET
	key.POST("/nws/stages", s.CreateNwsStage)        // CREATE
	key.PUT("/nws/stages/:nwsid", s.UpdateNwsStages) // UPDATE
}
