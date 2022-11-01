package locations

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (s Store) CreateLocations(c echo.Context) error {

	var nn LocationInfoCollection
	if err := c.Bind(&nn); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// ensure provider slug for all posted locations matches route param :provider
	// if not, return status unauthorized
	var routeProvider string
	if err := c.Bind(&routeProvider); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	for _, item := range nn.Items {
		if strings.ToUpper(routeProvider) != item.Provider {
			return c.String(
				http.StatusBadRequest,
				"location in post body has provider that does not match route param :provider",
			)
		}
	}

	cc, err := nn.LocationCollection() // convert to LocationCreatorCollection
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	ll, err := cc.Create(s.Connection) // run sql
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, ll)
}
