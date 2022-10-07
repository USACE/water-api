package visualizations

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

type Store struct {
	Connection *pgxpool.Pool
}

func Mount(conn *pgxpool.Pool, e *echo.Echo) {

	s := Store{Connection: conn} // database connection

	public := e.Group("")

	public.GET("/visualizations", s.ListVisualizations)
	public.GET("/visualizations/:visualization_slug", s.GetVisualization)
	public.POST("/visualizations", s.CreateVisualization)
	public.POST("/visualizations/:visualization_slug/assign", s.CreateOrUpdateVisualizationMapping)
}
