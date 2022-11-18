package charts

import (
	"log"
	"net/http"

	"github.com/USACE/water-api/api/app"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

type Store struct {
	Connection  *pgxpool.Pool
	ChartServer *ChartServer
}

func Mount(conn *pgxpool.Pool, e *echo.Echo, config *app.Config) {

	// d3-chart-server integration
	chartserver, err := NewChartServer(ChartServerConfig{URLString: config.ChartServerURL})
	if err != nil {
		log.Fatal(err.Error())
	}

	s := Store{Connection: conn, ChartServer: chartserver} // database connection

	public := e.Group("")

	// CHARTS; GLOBAL (ALL PROVIDERS) CONTEXT
	public.GET("/charts", s.ListCharts)            // LIST CHARTS
	public.GET("/charts/:chart", s.GetChartDetail) // GET CHART DETAILS (?format=svg renders chart)
	public.GET("/chart_types", func(c echo.Context) error { return c.JSON(http.StatusOK, s.ChartServer.Charts) })

	// CHARTS; SINGLE PROVIDER CONTEXT
	public.GET("/providers/:provider/charts", s.ListCharts)                                 // LIST CHARTS (Filter by Provider; Same as above with ?provider=<provider>)
	public.GET("/providers/:provider/charts/:chart", s.GetChartDetail)                      // GET CHART (filter by Provider)
	public.POST("/providers/:provider/charts", s.CreateCharts)                              // CREATE CHART(s)
	public.POST("/providers/:provider/charts/:chart/mapping", s.CreateOrUpdateChartMapping) // ADD CHART VARIABLE MAPPING
	public.DELETE("/providers/:provider/charts/:chart/mapping", s.DeleteChartMapping)       // DELETE CHART VARIABLE MAPPING
	public.DELETE("/providers/:provider/charts/:chart", s.DeleteChart)                      // DELETE CHART
}
