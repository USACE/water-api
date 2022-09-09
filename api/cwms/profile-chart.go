package cwms

import (
	"github.com/USACE/water-api/api/chartserver"
	"github.com/labstack/echo/v4"
)

func (s Store) GetProfileChart(c echo.Context) error {
	// TODO; Get location information needed to call chartserver
	//       mock in the data at this time.
	//
	input := chartserver.DamProfileChartInput{
		Pool:    15.0,
		Tail:    12.0,
		Inflow:  1200,
		Outflow: 600,
	}
	chart, err := s.ChartServer.GetDamProfileChart(input)
	if err != nil {
		return c.String(500, err.Error())
	}

	return c.String(200, chart)
}
