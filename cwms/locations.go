package cwms

import (
	"errors"
	"net/http"

	"github.com/USACE/water-api/cwms/models"
	"github.com/USACE/water-api/helpers"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
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
				return c.JSON(http.StatusNotFound, DefaultMessageNotFound)
			}
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, l)
	}
	// Otherwise Lookup By Slug
	l, err := models.GetLocationBySlug(s.Connection, &identifier)
	if err != nil {
		if pgxscan.NotFound(err) {
			return c.JSON(http.StatusNotFound, DefaultMessageNotFound)
		}
		return c.JSON(http.StatusInternalServerError, DefaultMessageInternalServerError)
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
					NewMessage("Locations not created. Location information conflicts with an existing location"))
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
		return c.JSON(http.StatusBadRequest, DefaultMessageBadRequest)
	}
	// Bind Payload (Updated Location Information)
	var uLocation models.Location
	if err := c.Bind(&uLocation); err != nil {
		return c.JSON(http.StatusBadRequest, DefaultMessageBadRequest)
	}
	// Compare Location ID from Route to Location ID from Payload
	if locationID != uLocation.ID {
		return c.JSON(
			http.StatusBadRequest,
			NewMessage("ID in Route Parameters does not match ID in payload"),
		)
	}
	l, err := models.UpdateLocation(s.Connection, &uLocation)
	if err != nil {
		if pgxscan.NotFound(err) {
			return c.JSON(http.StatusNotFound, DefaultMessageNotFound)
		}
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, l)

}

func (s Store) DeleteLocation(c echo.Context) error {
	locationID, err := uuid.Parse(c.Param("location_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, DefaultMessageBadRequest)
	}
	if err := models.DeleteLocation(s.Connection, &locationID); err != nil {
		return c.JSON(http.StatusInternalServerError, DefaultMessageInternalServerError)
	}
	return c.JSON(http.StatusOK, make(map[string]interface{}))
}
