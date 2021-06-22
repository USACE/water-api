package cwms

import (
	"net/http"

	"github.com/USACE/water-api/cwms/models"

	"github.com/labstack/echo/v4"
)

func (s Store) AssignStatesToLocations(c echo.Context) error {
	// @TODO: Query takes ~1m to run; This is too long to wait in context of web request; Investigate alternatives
	if err := models.AssignStatesToLocations(s.Connection); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, make(map[string]interface{}))
}
