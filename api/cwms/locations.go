package cwms

import (
	"errors"
	"net/http"

	"github.com/USACE/water-api/api/cwms/models"
	"github.com/USACE/water-api/api/helpers"
	"github.com/USACE/water-api/api/messages"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/labstack/echo/v4"
)

func (s Store) ListProjects(c echo.Context) error {
	pp, err := models.ListProjects(s.Connection)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, pp)
}

func (s Store) SearchLocations(c echo.Context) error {
	var f models.LocationFilter
	if err := c.Bind(&f); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if f.Q == nil || *f.Q == "" {
		return c.JSON(
			http.StatusBadRequest,
			messages.NewMessage("search string must be at one or more chacters long, provided in URL query parameter '?q='"),
		)
	}
	ll, err := models.SearchLocations(s.Connection, &f)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, messages.DefaultMessageInternalServerError)
	}
	return c.JSON(http.StatusOK, ll)
}

func (s Store) ListLocations(c echo.Context) error {
	// Get filters from query params kind_id= or office_id=
	var f models.LocationFilter
	if err := c.Bind(&f); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	ll, err := models.ListLocations(s.Connection, &f)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ll)
}

func (s Store) GetLocation(c echo.Context) error {
	identifier := c.Param("location_id")
	// Lookup By Location UUID if location_id parseable as UUID
	if locationID, err := uuid.Parse(identifier); err == nil {
		l, err := models.GetLocationByID(s.Connection, &locationID)
		if err != nil {
			if pgxscan.NotFound(err) {
				return c.JSON(http.StatusNotFound, messages.DefaultMessageNotFound)
			}
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, l)
	}
	// Otherwise Lookup By Slug
	l, err := models.GetLocationBySlug(s.Connection, &identifier)
	if err != nil {
		if pgxscan.NotFound(err) {
			return c.JSON(http.StatusNotFound, messages.DefaultMessageNotFound)
		}
		return c.JSON(http.StatusInternalServerError, messages.DefaultMessageInternalServerError)
	}
	return c.JSON(http.StatusOK, l)
}

func (s Store) CreateLocations(c echo.Context) error {
	var lc models.LocationCollection
	if err := c.Bind(&lc); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// Assign Unique Slugs
	for idx := range lc.Items {
		_s, err := helpers.NextUniqueSlug(s.Connection, "location", "slug", lc.Items[idx].Name, "", "")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		lc.Items[idx].Slug = _s
	}
	// Create Locations
	ll, err := models.CreateLocations(s.Connection, lc)
	if err != nil {
		// If Error was postgres error, return error message based on error code
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return c.JSON(
					http.StatusBadRequest,
					messages.NewMessage("Locations not created. Location information conflicts with an existing location"))
			}
		}
		// If not explicit error, return string of error message for debugging
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, ll)
}

// CreateLocationsByOffice
func (s Store) CreateLocationsByOffice(c echo.Context) error {
	office_symbol := c.Param("office_symbol")
	var lc models.LocationCollection
	if err := c.Bind(&lc); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// Assign Unique Slugs
	for idx := range lc.Items {
		_s, err := helpers.NextUniqueSlug(s.Connection, "location", "slug", lc.Items[idx].Name, "", "")
		if err != nil {
			lc.Items[idx].Slug = _s
		} else {
			lc.Items[idx].Slug = slug.Make(lc.Items[idx].Name)
		}
	}
	ll, err := models.CreateLocationsByOffice(s.Connection, lc, &office_symbol)
	if err != nil {
		// If Error was postgres error, return error message based on error code
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return c.JSON(
					http.StatusBadRequest,
					messages.NewMessage("Locations not created. Location information conflicts with an existing location"))
			}
		}
		// If not explicit error, return string of error message for debugging
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, ll)
}

func (s Store) UpdateLocation(c echo.Context) error {
	// Location ID From Route
	locationID, err := uuid.Parse(c.Param("location_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, messages.DefaultMessageBadRequest)
	}
	// Bind Payload (Updated Location Information)
	var uLocation models.Location
	if err := c.Bind(&uLocation); err != nil {
		return c.JSON(http.StatusBadRequest, messages.DefaultMessageBadRequest)
	}
	// Compare Location ID from Route to Location ID from Payload
	if locationID != uLocation.ID {
		return c.JSON(
			http.StatusBadRequest,
			messages.NewMessage("ID in Route Parameters does not match ID in payload"),
		)
	}
	l, err := models.UpdateLocation(s.Connection, &uLocation)
	if err != nil {
		if pgxscan.NotFound(err) {
			return c.JSON(http.StatusNotFound, messages.DefaultMessageNotFound)
		}
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, l)

}

// UpdateLocationByOffice
func (s Store) UpdateLocationByOffice(c echo.Context) error {
	office_symbol := c.Param("office_symbol")
	location_slug := c.Param("location_slug")
	var uLocation models.Location
	if err := c.Bind(&uLocation); err != nil {
		return c.JSON(http.StatusBadRequest, messages.DefaultMessageBadRequest)
	}
	// Compare slug from Route to Location Slug from Payload
	if location_slug != uLocation.Slug {
		return c.JSON(
			http.StatusBadRequest,
			messages.NewMessage("Slug in Route Parameters does not match slug in payload"),
		)
	}
	l, err := models.UpdateLocationByOffice(s.Connection, &uLocation, &office_symbol)
	if err != nil {
		if pgxscan.NotFound(err) {
			return c.JSON(http.StatusNotFound, messages.DefaultMessageNotFound)
		}
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, l)
}

func (s Store) SyncLocations(c echo.Context) error {
	var lc models.LocationCollection
	if err := c.Bind(&lc); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// Assign Unique Slugs
	for idx := range lc.Items {
		_s, err := helpers.NextUniqueSlug(s.Connection, "location", "slug", lc.Items[idx].Name, "", "")
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		lc.Items[idx].Slug = _s
	}
	sl, err := models.SyncLocations(s.Connection, lc)
	if err != nil {
		// The server error results from UPDATE; try to create
		// Assign Unique Slugs
		for idx := range lc.Items {
			_s, err := helpers.NextUniqueSlug(s.Connection, "location", "slug", lc.Items[idx].Name, "", "")
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
			lc.Items[idx].Slug = _s
		}
		cl, err := models.CreateLocations(s.Connection, lc)
		// If Error was postgres error, return error message based on error code
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return c.JSON(
					http.StatusBadRequest,
					messages.NewMessage(pgErr.Error()))
			}
		}
		return c.JSON(http.StatusCreated, cl)
	}
	return c.JSON(http.StatusAccepted, sl)
}

func (s Store) DeleteLocation(c echo.Context) error {
	locationID, err := uuid.Parse(c.Param("location_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, messages.DefaultMessageBadRequest)
	}
	if err := models.DeleteLocation(s.Connection, &locationID); err != nil {
		return c.JSON(http.StatusInternalServerError, messages.DefaultMessageInternalServerError)
	}
	return c.JSON(http.StatusOK, make(map[string]interface{}))
}

// DeleteLocationByOffice
func (s Store) DeleteLocationByOffice(c echo.Context) error {
	// Need to parse and input parameters
	location_slug := c.Param("location_slug")
	office_symbol := c.Param("office_symbol")
	if err := models.DeleteLocationByOffice(s.Connection, location_slug, office_symbol); err != nil {
		return c.JSON(http.StatusInternalServerError, messages.DefaultMessageInternalServerError)
	}
	return c.JSON(http.StatusOK, make(map[string]interface{}))
}
