package usgs

import (
	"net/http"

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
