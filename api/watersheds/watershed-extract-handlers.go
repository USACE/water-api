package watersheds

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/USACE/water-api/api/timeseries"
	"github.com/USACE/water-api/api/watersheds/models"
)

// WatershedExtract
func (s Store) WatershedExtract(c echo.Context) error {
	wslug := c.Param("watershed_slug")

	// Time Window
	ab := map[string]string{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, &ab); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	tw, err := timeseries.CreateTimeWindow(ab["after"], ab["before"])

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	rows, err := models.WatershedExtract(s.Connection, wslug, &tw)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	enc := json.NewEncoder(c.Response())
	for _, row := range rows {
		if err := enc.Encode(row); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		c.Response().Flush()
	}
	return c.NoContent(http.StatusOK)
}
