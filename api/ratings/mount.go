package ratings

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

	// RATING CRUD OPERATIONS
	public.GET("/ratings/:slug/:method/:ival", s.Rate1Var)

}
