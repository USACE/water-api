package timeseries

import (
	"net/http"
	"strings"

	"github.com/USACE/water-api/api/messages"
	"github.com/USACE/water-api/api/timeseries/models"
	"github.com/labstack/echo/v4"
)

func (s Store) ListTimeseries(c echo.Context) error {
	// Get filters from query provider= or datatype=
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

	_, err := models.CreateTimeseries(s.Connection, tsc)

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

	return c.NoContent(http.StatusAccepted)
}
