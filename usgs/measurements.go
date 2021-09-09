package usgs

import (
	"net/http"
	"time"

	"github.com/USACE/water-api/timeseries"
	"github.com/USACE/water-api/usgs/models"
	"github.com/labstack/echo/v4"
)

// CreateOrUpdateTimeseriesMeasurements
func (s Store) CreateOrUpdateMeasurements(c echo.Context) error {
	var pm models.ParameterMeasurementCollection
	if err := c.Bind(&pm); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ss, err := models.CreateOrUpdateMeasurements(s.Connection, pm)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusAccepted, ss)
}

// ListMeasurements
func (s Store) ListUSGSMeasurements(c echo.Context) error {
	site_number := c.Param("site_number")
	parameters := c.QueryParams()["parameter"]

	// Time Window
	var tw timeseries.TimeWindow
	a, b := c.QueryParam("after"), c.QueryParam("before")
	// If after or before are not provided return last 7 days of data from current time
	if a == "" || b == "" {
		tw.Before = time.Now()
		tw.After = tw.Before.AddDate(0, 0, -7)
	} else {
		// Attempt to parse query param "after"
		tA, err := time.Parse(time.RFC3339, a)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		tw.After = tA
		// Attempt to parse query param "before"
		tB, err := time.Parse(time.RFC3339, b)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		tw.Before = tB
	}
	mc, err := models.ListUSGSMeasurements(s.Connection, &site_number, parameters, &tw)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusAccepted, mc)
}
