package cwms

import (
	"net/http"

	"github.com/USACE/water-api/api/cwms/models"
	"github.com/labstack/echo/v4"
)

func (s Store) SyncTimeseries(c echo.Context) error {
	var tsc models.TimeseriesCollection
	if err := c.Bind(&tsc); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	s_tsc, err := models.SyncTimeseries(s.Connection, tsc)

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusAccepted, s_tsc)
}
