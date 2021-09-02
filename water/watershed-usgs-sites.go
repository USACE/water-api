package water

import (
	"net/http"
	"strings"

	"github.com/USACE/water-api/water/models"

	"github.com/labstack/echo/v4"
)

func (s Store) CreateWatershedSiteParameter(c echo.Context) error {

	var wsp models.WatershedSiteParameter
	if err := c.Bind(&wsp); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err := models.CreateWatershedSiteParameter(s.Connection, &wsp)
	if err != nil {
		if strings.Contains(string(err.Error()), "duplicate key value") {
			// return 422
			return c.JSON(http.StatusUnprocessableEntity, make(map[string]interface{}))
		}
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, make(map[string]interface{}))
}

func (s Store) DeleteWatershedSiteParameter(c echo.Context) error {

	var wsp models.WatershedSiteParameter
	if err := c.Bind(&wsp); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err := models.DeleteWatershedSiteParameter(s.Connection, &wsp)
	if err != nil {
		// if strings.Contains(string(err.Error()), "duplicate key value") {
		// 	// return 422
		// 	return c.String(http.StatusUnprocessableEntity, err.Error())
		// }
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, make(map[string]interface{}))
}
