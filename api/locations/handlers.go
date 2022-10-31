package locations

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s Store) CreateLocations(c echo.Context) error {

	var lc LocationCollection
	if err := c.Bind(&lc); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	cc, err := lc.LocationCreatorCollection() // convert to LocationCreatorCollection
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	ll, err := cc.Create(s.Connection) // run sql
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, ll)
}
