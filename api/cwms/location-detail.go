package cwms

import (
	"net/http"

	"github.com/USACE/water-api/api/cwms/models"
	"github.com/USACE/water-api/api/messages"
	"github.com/georgysavva/scany/pgxscan"

	"github.com/labstack/echo/v4"
)

func (s Store) GetLocationDetail(c echo.Context) error {
	identifier := c.Param("location_slug")
	// Lookup By Location UUID if location_id parseable as UUID
	l, err := models.GetLocationDetail(s.Connection, &identifier)
	if err != nil {
		if pgxscan.NotFound(err) {
			return c.JSON(http.StatusNotFound, messages.DefaultMessageNotFound)
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, l)
}
