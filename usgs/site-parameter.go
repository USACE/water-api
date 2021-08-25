package usgs

import (
	"net/http"
	"strings"

	"github.com/USACE/water-api/usgs/models"

	"github.com/labstack/echo/v4"
)

func (s Store) CreateSiteParameters(c echo.Context) error {

	var sp models.SiteParameterCollection
	if err := c.Bind(&sp); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	ss, err := models.CreateSiteParameters(s.Connection, sp.Items)
	if err != nil {
		if strings.Contains(string(err.Error()), "duplicate key value") {
			// return 422
			return c.String(http.StatusUnprocessableEntity, err.Error())
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, ss)
}
