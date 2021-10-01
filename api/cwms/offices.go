package cwms

import (
	"net/http"

	"github.com/USACE/water-api/cwms/models"

	"github.com/labstack/echo/v4"
)

func (s Store) ListOffices(c echo.Context) error {
	oo, err := models.ListOffices(s.Connection)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, oo)
}
