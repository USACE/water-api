package timeseries

import (
	"net/http"
	"strings"

	"github.com/USACE/water-api/api/messages"
	"github.com/USACE/water-api/api/timeseries/models"
	"github.com/labstack/echo/v4"
)

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
			return c.JSON(http.StatusBadRequest, messages.NewMessage("measurements not provided in proper format"))
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusAccepted)

}
