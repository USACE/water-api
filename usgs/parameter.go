package usgs

import (
	"net/http"

	"github.com/USACE/water-api/usgs/models"

	"github.com/labstack/echo/v4"
)

func (s Store) ListParameters(c echo.Context) error {
	pp, err := models.ListParameters(s.Connection)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, pp)
}

// NOT USED - SAVE FOR WATERSHED/SITE/PARAM Enabled
// func (s Store) ListParametersEnabled(c echo.Context) error {
// 	pp, err := models.ListParametersEnabled(s.Connection)
// 	if err != nil {
// 		return c.String(http.StatusInternalServerError, err.Error())
// 	}
// 	return c.JSON(http.StatusOK, pp)
// }
