package cwms

import (
	"net/http"

	"github.com/USACE/water-api/cwms/models"

	"github.com/labstack/echo/v4"
)

func (s Store) ListSites(c echo.Context) error {

	// Get filter from query params state_id =
	var sf models.SiteFilter
	if err := c.Bind(&sf); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	ss, err := models.ListSites(s.Connection, &sf)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ss)
}
