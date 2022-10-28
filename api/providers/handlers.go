package providers

import (
	"net/http"

	"github.com/USACE/water-api/api/providers/models"
	"github.com/labstack/echo/v4"
)

func (s Store) ListProviders(c echo.Context) error {
	pp, err := models.ListProviders(s.Connection)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, pp)
}

// ListDatasources
func (s Store) ListDatasources(c echo.Context) error {
	p := c.QueryParam("provider")
	d := c.QueryParam("datasource_type")
	dd, err := models.ListDatasources(s.Connection, p, d)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dd)
}
