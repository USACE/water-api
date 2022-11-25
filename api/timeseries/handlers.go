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

	// ensure provider slug for all posted locations matches route param :provider
	// if not, return status unauthorized
	routeProvider := c.Param("provider")
	for _, item := range tsc.Items {
		if !strings.EqualFold(routeProvider, item.Provider) {
			return c.String(
				http.StatusBadRequest,
				"timeseries in post body has provider that does not match route param :provider",
			)
		}
	}

	tt, err := tsc.Create(s.Connection, strings.ToLower(routeProvider))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, messages.DefaultMessageInternalServerError)
	}

	// If 0 new timeseries were created, return a RESTful 200
	if len(tt) == 0 {
		return c.JSON(http.StatusOK, tt)
	}

	// If at least 1 timeseries was created, return 201 with array of new timeseries
	return c.JSON(http.StatusCreated, tt)
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

	tt, err := tsc.CreateOrUpdateTimeseriesMeasurements(s.Connection)
	if err != nil {
		if strings.Contains(err.Error(), "no records updated") {
			return c.JSON(http.StatusBadRequest, messages.NewMessage("measurements not provided in proper format"))
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusAccepted, tt)

}

func (s Store) UpdateTimeseries(c echo.Context) error {
	var tsc models.TimeseriesCollection
	if err := c.Bind(&tsc); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// ensure provider slug for all posted locations matches route param :provider
	// if not, return status unauthorized
	routeProvider := c.Param("provider")
	for _, item := range tsc.Items {
		if !strings.EqualFold(routeProvider, item.Provider) {
			return c.String(
				http.StatusBadRequest,
				"timeseries in post body has provider that does not match route param :provider",
			)
		}
	}

	tt, err := tsc.Update(s.Connection)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, messages.DefaultMessageInternalServerError)
	}

	return c.JSON(http.StatusOK, tt)
}

func (s Store) DeleteTimeseries(c echo.Context) error {
	var tsc models.TimeseriesCollection
	if err := c.Bind(&tsc); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// ensure provider slug for all posted locations matches route param :provider
	// if not, return status unauthorized
	routeProvider := c.Param("provider")
	for _, item := range tsc.Items {
		if !strings.EqualFold(routeProvider, item.Provider) {
			return c.String(
				http.StatusBadRequest,
				"timeseries in post body has provider that does not match route param :provider",
			)
		}
	}

	if err := tsc.Delete(s.Connection); err != nil {
		return c.JSON(http.StatusInternalServerError, messages.DefaultMessageInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{})
}
