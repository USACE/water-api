package usgs

import (
	"net/http"
	"strings"

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
	return c.JSON(http.StatusCreated, ss)
}

// ListMeasurements
func (s Store) ListUSGSMeasurements(c echo.Context) error {
	site_number := c.Param("site_number")
	parameters := make([]string, 0)
	for _, element := range c.QueryParams()["parameter"] {
		// Split and trim
		s := strings.Split(element, ",")
		for i, e := range s {
			s[i] = strings.TrimSpace(e)
		}
		parameters = append(parameters, s...)
	}
	// Remove duplicates
	allKeys := make(map[string]bool)
	filter_list := []string{}
	for _, item := range parameters {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			filter_list = append(filter_list, item)
		}
	}
	parameters = filter_list

	// Time Window
	a, b := c.QueryParam("after"), c.QueryParam("before")
	tw, err := timeseries.CreateTimeWindow(a, b)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	mc, err := models.ListUSGSMeasurements2(s.Connection, &site_number, parameters, &tw)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, mc)
}
