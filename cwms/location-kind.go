package cwms

import (
	"net/http"

	"github.com/USACE/water-api/cwms/models"

	"github.com/labstack/echo/v4"
)

func (s Store) ListLocationKind(c echo.Context) error {
	lk, err := models.ListLocationKind(s.Connection)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, lk)
}
