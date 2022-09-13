package cwms

import (
	"net/http"

	"github.com/USACE/water-api/api/cwms/models"
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

func (s Store) CreateOrUpdateTimeseries(c echo.Context) error {
	var tsc models.TimeseriesCollection
	if err := c.Bind(&tsc); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	s_tsc, err := models.CreateOrUpdateTimeseries(s.Connection, tsc)

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusAccepted, s_tsc)
}
