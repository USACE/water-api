package locations

import (
	"errors"
	"net/http"

	"github.com/USACE/water-api/api/cwms/models"
	"github.com/USACE/water-api/api/helpers"
	"github.com/USACE/water-api/api/messages"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/labstack/echo/v4"
)

func (s Store) ListLocations(c echo.Context) error {
	// Get filters from query params
	var f models.LocationFilter
	if err := c.Bind(&f); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// Search Query must not be blank
	if f.Q != nil && *f.Q == "" {
		return c.JSON(
			http.StatusBadRequest,
			messages.NewMessage("search string must be at one or more chacters long, provided in URL query parameter '?q='"),
		)
	}
	ll, err := models.ListLocations(s.Connection, &f)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ll)
}

func (s Store) GetLocation(c echo.Context) error {
	// Get filters from query params
	var f models.LocationFilter
	if err := c.Bind(&f); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	l, err := models.GetLocation(s.Connection, &f)
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
