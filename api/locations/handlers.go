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
	routeProvider := c.Param("provider")
	for _, item := range nn.Items {
		if !strings.EqualFold(routeProvider, item.Provider) {
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

	// If 0 new locations were created, return a RESTful 200
	if len(ll) == 0 {
		return c.JSON(http.StatusOK, ll)
	}
	// If at least 1 location was created, return 201 with array of new locations
	return c.JSON(http.StatusCreated, ll)
}

func (s Store) UpdateLocations(c echo.Context) error {

	var nn LocationInfoCollection
	if err := c.Bind(&nn); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// ensure provider slug for all posted locations matches route param :provider
	// if not, return status unauthorized
	routeProvider := c.Param("provider")
	for _, item := range nn.Items {
		if !strings.EqualFold(routeProvider, item.Provider) {
			return c.String(
				http.StatusBadRequest,
				"location in post body has provider that does not match route param :provider",
			)
		}
	}

	cc, err := nn.LocationCollection() // convert to LocationCollection
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	ll, err := cc.Update(s.Connection) // run sql
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, ll)
}

func (s Store) ListLocations(c echo.Context) error {
	// Get filters from query params kind_id= or office_id=
	var f LocationFilter
	if err := c.Bind(&f); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	ll, err := ListLocations(s.Connection, &f)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if len(ll) == 0 {
		return c.JSON(http.StatusNotFound, ll)
	}

	return c.JSON(http.StatusOK, ll)
}

func (s Store) GetLocation(c echo.Context) error {
	// Get filters from query params; The :location route parameter is all that is needed,
	// as this is globally unique for a location.  Binding LocationFilter for shared behavior
	// with ListLocationsQuery()
	var f LocationFilter
	if err := c.Bind(&f); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	l, err := GetLocation(s.Connection, &f)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, l)
}
