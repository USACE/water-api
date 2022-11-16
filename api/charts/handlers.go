package charts

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/USACE/water-api/api/messages"
	"github.com/labstack/echo/v4"
)

func (s Store) ListCharts(c echo.Context) error {

	var f ChartFilter
	if err := c.Bind(&f); err != nil {
		return c.JSON(http.StatusBadRequest, messages.NewMessage(err.Error()))
	}

	tt, err := ListCharts(s.Connection, &f)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, tt)
}

func (s Store) GetChartDetail(c echo.Context) error {

	var f ChartFilter
	if err := c.Bind(&f); err != nil {
		return c.JSON(http.StatusBadRequest, messages.NewMessage(err.Error()))
	}
	t, err := GetChartDetail(s.Connection, &f)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.JSON(http.StatusNotFound, messages.DefaultMessageNotFound)
		}
		return c.JSON(http.StatusInternalServerError, messages.NewMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, t)
}

// CreateChart creates one or more new charts
func (s Store) CreateCharts(c echo.Context) error {

	var cc ChartCollection
	if err := c.Bind(&cc); err != nil {
		return c.JSON(http.StatusBadRequest, messages.NewMessage(err.Error()))
	}

	// Validate incoming payload. Check:
	//   1. Verify provider in payload body matches route param :provider
	//   2. "type" of chart in body is supported by backend ChartServer
	routeProvider := c.Param("provider")
	for _, item := range cc.Items {
		// 1. Verify provider in payload body matches route param :provider
		if !strings.EqualFold(routeProvider, item.Provider) {
			return c.JSON(
				http.StatusBadRequest,
				messages.NewMessage(
					fmt.Sprintf(
						"chart in post body has provider (%s) that does not match route param :provider (%s)",
						item.Provider,
						routeProvider,
					),
				),
			)
		}
		// 2. "type" of chart in body is supported by backend ChartServer
		if _, ok := s.ChartServer.ChartMap[strings.ToLower(item.Type)]; !ok {
			return c.JSON(http.StatusBadRequest, messages.NewMessage(fmt.Sprintf("unsupported chart type: %s", item.Type)))
		}
	}

	new, err := cc.Create(s.Connection)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, messages.NewMessage(err.Error()))
	}
	return c.JSON(http.StatusCreated, &new)
}

func (s Store) CreateOrUpdateChartMapping(c echo.Context) error {

	provider, chart := c.Param("provider"), c.Param("chart")

	var mc ChartMappingCollection
	if err := c.Bind(&mc); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	err := CreateOrUpdateChartMapping(s.Connection, &provider, &chart, &mc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, messages.NewMessage(err.Error()))
	}
	return c.JSON(http.StatusAccepted, messages.DefaultMessageOK)
}

func (s Store) DeleteChartMapping(c echo.Context) error {

	provider, chart := c.Param("provider"), c.Param("chart")

	var mc ChartMappingCollection
	if err := c.Bind(&mc); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	err := DeleteChartMapping(s.Connection, &provider, &chart, &mc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, messages.NewMessage(err.Error()))
	}
	return c.JSON(http.StatusAccepted, messages.DefaultMessageOK)
}

func (s Store) DeleteChart(c echo.Context) error {

	provider, chart := c.Param("provider"), c.Param("chart")

	if err := DeleteChart(s.Connection, &provider, &chart); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, messages.DefaultMessageOK)
}
