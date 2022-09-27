package cwms

import (
	"net/http"
	"strings"

	"github.com/USACE/water-api/api/cwms/models"
	"github.com/USACE/water-api/api/messages"
	"github.com/labstack/echo/v4"
)

func (s Store) ListTimeseries(c echo.Context) error {
	// Get filters from query provider= or datasource_type=
	var f models.TimeseriesFilter
	if err := c.Bind(&f); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	ll, err := models.ListTimeseries(s.Connection, &f)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ll)
}

func (s Store) CreateTimeseries(c echo.Context) error {
	var tsc models.TimeseriesCollection
	if err := c.Bind(&tsc); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	s_tsc, err := models.CreateTimeseries(s.Connection, tsc)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return c.JSON(http.StatusBadRequest, messages.NewMessage("duplicate timeseries submitted"))
		}
		if strings.Contains(err.Error(), "constraint") {
			return c.JSON(http.StatusBadRequest, messages.DefaultMessageBadRequest)
		}
		//return c.String(http.StatusInternalServerError, err.Error())
		return c.JSON(http.StatusInternalServerError, messages.DefaultMessageInternalServerError)
	}

	return c.JSON(http.StatusCreated, s_tsc)
}

func (s Store) CreateOrUpdateTimeseriesMeasurements(c echo.Context) error {
	var tsc models.TimeseriesCollection
	if err := c.Bind(&tsc); err != nil {
		if strings.Contains(err.Error(), "parsing time") {
			return c.JSON(http.StatusBadRequest, messages.NewMessage("incorrect time format, use YYYY-MM-DDTHH:MM:SS-HH:MM"))
		}
		if strings.Contains(err.Error(), "measurements.values") {
			return c.JSON(http.StatusBadRequest, messages.NewMessage("incorrect value type, values should be an array of integers (99) or floats (99.99)"))
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	_, err := models.CreateOrUpdateTimeseriesMeasurements(s.Connection, tsc)

	if err != nil {
		if strings.Contains(err.Error(), "no records updated") {
			return c.JSON(http.StatusBadRequest, messages.NewMessage("measurements no provided in proper format"))
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusAccepted)

}
