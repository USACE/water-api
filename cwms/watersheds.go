package cwms

import (
	"net/http"

	"github.com/USACE/water-api/cwms/models"
	"github.com/USACE/water-api/messages"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
)

// ListWatersheds returns an array of Watersheds
func (s Store) ListWatersheds(c echo.Context) error {
	ww, err := models.ListWatersheds(s.Connection)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ww)

}

// GetWatershed returns a single Watershed
func (s Store) GetWatershed(c echo.Context) error {
	id, err := uuid.Parse(c.Param("watershed_id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	w, err := models.GetWatershed(s.Connection, &id)
	if err != nil {
		if pgxscan.NotFound(err) {
			return c.JSON(http.StatusNotFound, messages.DefaultMessageNotFound)
		}
		return c.JSON(http.StatusInternalServerError, messages.DefaultMessageInternalServerError)
	}
	return c.JSON(http.StatusOK, w)
}

// CreateWatershed creates a new watershed
func (s Store) CreateWatershed(c echo.Context) error {
	var w models.Watershed
	if err := c.Bind(&w); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	newWatershed, err := models.CreateWatershed(s.Connection, &w)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, newWatershed)
}

// UpdateWatershed creates a new watershed
func (s Store) UpdateWatershed(c echo.Context) error {
	// Watershed Slug from route params
	wID, err := uuid.Parse(c.Param("watershed_id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// Payload
	var w models.Watershed
	if err := c.Bind(&w); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// Check route params v. payload
	if wID != w.ID {
		return c.String(http.StatusBadRequest, "watershed_id in URL does not match request body")
	}
	wUpdated, err := models.UpdateWatershed(s.Connection, &w)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, wUpdated)

}

// DeleteWatershed creates a new watershed
func (s Store) DeleteWatershed(c echo.Context) error {
	wID, err := uuid.Parse(c.Param("watershed_id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = models.DeleteWatershed(s.Connection, &wID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, make(map[string]interface{}))
}

// UndeleteWatershed restores a deleted watershed
func (s Store) UndeleteWatershed(c echo.Context) error {
	wID, err := uuid.Parse(c.Param("watershed_id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	w, err := models.UndeleteWatershed(s.Connection, &wID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, w)

}
