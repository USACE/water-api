package water

import (
	"net/http"

	"github.com/USACE/water-api/api/water/models"
	"github.com/labstack/echo/v4"
)

func (s Store) ListProviders(c echo.Context) error {

	pp, err := models.ListProviders(s.Connection)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, pp)
}
