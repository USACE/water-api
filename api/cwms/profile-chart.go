package cwms

import (
	"net/http"
	"strings"

	"github.com/USACE/water-api/api/chartserver"
	"github.com/labstack/echo/v4"
)

func (s Store) GetProfileChart(c echo.Context) error {
	// TODO; Get location information needed to call chartserver
	//       mock in the data at this time.
	//

	format := c.QueryParam("format")

	locationSlug := c.Param("location_slug")

	//v, err := water.GetVisualizationByLocation(s.Connection, &locationSlug, &visualizationTypeId)
	v, err := chartserver.GetDamProfileByLocation(s.Connection, &locationSlug)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	input := chartserver.DamProfileChartInput{
		Pool:      v.Pool,
		Tail:      v.Tail,
		Inflow:    v.Inflow,
		Outflow:   v.Outflow,
		DamTop:    v.DamTop,
		DamBottom: v.DamBottom,
		Levels:    v.Levels,
	}
	chart, err := s.ChartServer.DamProfileChart(input)
	if err != nil {
		return c.String(500, err.Error())
	}

	// return json if querystring specified
	if strings.ToLower(format) == "json" {
		return c.JSON(http.StatusOK, v)
	}
	// otherwise display the chart
	return c.HTML(200, chart)
}
