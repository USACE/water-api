package charts

import (
	"net/http"
	"strings"

	"github.com/USACE/water-api/api/charts/models"
	"github.com/USACE/water-api/api/messages"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/labstack/echo/v4"
)

func (s Store) ListCharts(c echo.Context) error {

	vv, err := models.ListCharts(s.Connection)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, vv)
}

func (s Store) GetChart(c echo.Context) error {

	vSlug := c.Param("chart_slug")
	t, err := models.GetChart(s.Connection, &vSlug)

	if err != nil {
		if pgxscan.NotFound(err) {
			return c.JSON(http.StatusNotFound, messages.DefaultMessageNotFound)
		}
		// return c.JSON(http.StatusInternalServerError, messages.DefaultMessageInternalServerError)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, t)
}

// CreateChart creates a single new Chart
func (s Store) CreateChart(c echo.Context) error {

	var v models.Chart
	if err := c.Bind(&v); err != nil {
		return c.JSON(http.StatusBadRequest, messages.DefaultMessageBadRequest)
	}
	vNew, err := models.CreateChart(s.Connection, &v)
	if err != nil {
		if strings.Contains(err.Error(), `null value in column "location_id"`) {
			return c.JSON(http.StatusBadRequest, messages.NewMessage("invalid location_slug"))
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, &vNew)
}

func (s Store) CreateOrUpdateChartMapping(c echo.Context) error {

	slug := c.Param("chart_slug")

	var vm models.ChartMappingCollection
	if err := c.Bind(&vm); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	_, err := models.CreateOrUpdateChartMapping(s.Connection, vm, &slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusAccepted)
}

func (s Store) DeleteChart(c echo.Context) error {
	vSlug := c.Param("chart_slug")

	err := models.DeleteChart(s.Connection, &vSlug)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, make(map[string]interface{}))
}
