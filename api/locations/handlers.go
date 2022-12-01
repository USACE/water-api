package locations

import (
	"net/http"
	"strings"

	"github.com/USACE/water-api/api/messages"
	"github.com/labstack/echo/v4"
)

func (s Store) CreateLocations(c echo.Context) error {

	var nn LocationInfoCollection
	if err := c.Bind(&nn); err != nil {
		return c.JSON(http.StatusBadRequest, messages.NewMessage(err.Error()))
	}

	// ensure provider slug for all posted locations matches route param :provider
	// if not, return status unauthorized
	routeProvider := c.Param("provider")
	for _, item := range nn.Items {
		if !strings.EqualFold(routeProvider, item.Provider) {
			return c.JSON(
				http.StatusBadRequest,
				"location in post body has provider that does not match route param :provider",
			)
		}
	}

	cc, err := nn.LocationCollection() // convert to LocationCollection
	if err != nil {
		return c.JSON(http.StatusInternalServerError, messages.NewMessage(err.Error()))
	}

	ll, err := cc.Create(s.Connection) // run sql
	if err != nil {
		return c.JSON(http.StatusInternalServerError, messages.NewMessage(err.Error()))
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
		return c.JSON(http.StatusBadRequest, messages.NewMessage(err.Error()))
	}

	// ensure provider slug for all posted locations matches route param :provider
	// if not, return status unauthorized
	routeProvider := c.Param("provider")
	for _, item := range nn.Items {
		if !strings.EqualFold(routeProvider, item.Provider) {
			return c.JSON(
				http.StatusBadRequest,
				"location in post body has provider that does not match route param :provider",
			)
		}
	}

	cc, err := nn.LocationCollection() // convert to LocationCollection
	if err != nil {
		return c.JSON(http.StatusInternalServerError, messages.NewMessage(err.Error()))
	}

	ll, err := cc.Update(s.Connection) // run sql
	if err != nil {
		return c.JSON(http.StatusInternalServerError, messages.NewMessage(err.Error()))
	}

	return c.JSON(http.StatusOK, ll)
}

func (s Store) ListLocations(c echo.Context) error {
	// Get filters from query params kind_id= or office_id=
	var f LocationFilter
	if err := c.Bind(&f); err != nil {
		return c.JSON(http.StatusBadRequest, messages.NewMessage(err.Error()))
	}

	ll, err := ListLocations(s.Connection, &f)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, messages.NewMessage(err.Error()))
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
		return c.JSON(http.StatusBadRequest, messages.NewMessage(err.Error()))
	}

	l, err := GetLocation(s.Connection, &f)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, messages.NewMessage(err.Error()))
	}

	return c.JSON(http.StatusOK, l)
}

func (s Store) DeleteLocations(c echo.Context) error {

	var ic LocationInfoCollection
	if err := c.Bind(&ic); err != nil {
		return c.JSON(http.StatusBadRequest, messages.NewMessage(err.Error()))
	}

	// ensure provider slug for all posted locations matches route param :provider
	// if not, return status unauthorized
	routeProvider := c.Param("provider")
	for _, item := range ic.Items {
		if !strings.EqualFold(routeProvider, item.Provider) {
			return c.JSON(
				http.StatusBadRequest,
				"location in post body has provider that does not match route param :provider",
			)
		}
	}

	// Delete
	if err := Delete(s.Connection, ic); err != nil {
		return c.JSON(http.StatusInternalServerError, messages.NewMessage(err.Error()))
	}

	return c.JSON(http.StatusOK, messages.DefaultMessageOK)
}

// DeleteOneLoation deletes a single location using the routing parameter
// :location, corresponding to a location's globally unique slug.
func (s Store) DeleteOneLocation(c echo.Context) error {

	provider, slug := c.Param("provider"), c.Param("location")
	if provider == "" || slug == "" {
		return c.JSON(http.StatusBadRequest, messages.DefaultMessageBadRequest)
	}

	// Create Equivalent LocationInfoCollection expected by Delete func
	ic := LocationInfoCollection{
		Items: []LocationInfo{
			{Slug: slug, Provider: provider},
		},
	}

	// Delete
	if err := Delete(s.Connection, ic); err != nil {
		return c.JSON(http.StatusInternalServerError, messages.NewMessage(err.Error()))
	}

	return c.JSON(http.StatusOK, messages.DefaultMessageOK)
}
