package charts

import (
	"github.com/USACE/water-api/api/app"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

type Store struct {
	Connection *pgxpool.Pool
}

func Mount(conn *pgxpool.Pool, e *echo.Echo, config *app.Config) {

	s := Store{Connection: conn} // database connection

	public := e.Group("")

	public.GET("/charts", s.ListCharts)
	public.GET("/charts/:chart_slug", s.GetChart)
	public.POST("/charts", s.CreateChart)
	public.POST("/charts/:chart_slug/assign", s.CreateOrUpdateChartMapping)
	public.DELETE("/charts/:chart_slug", s.DeleteChart)
}
