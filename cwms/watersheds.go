package cwms

import (
	"net/http"

	"github.com/USACE/water-api/cwms/models"
	"github.com/USACE/water-api/messages"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/paulmach/orb/geojson"

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
func (s Store) UpdateWatershedGeometry(c echo.Context) error {
	var (
		err error
		id  uuid.UUID
		ws  *models.Watershed
		wf  geojson.Feature
	)
	p1 := c.ParamNames()[0] // Just incase the parameter name is changed for the watershed id
	if id, err = uuid.Parse(c.Param(p1)); err != nil {
		return c.JSON(http.StatusBadRequest, "UUID parsing error: "+err.Error())
	}
	// var wf geojson.Feature
	if err = c.Bind(&wf); err != nil {
		return c.JSON(http.StatusBadRequest, "Binding error: "+err.Error())
	}
	if ws, err = models.UpdateWatershedGeometry(s.Connection, &id, &wf); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, ws)
}

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

// UploadWatersheds handler for models.UploadWatersheds
// slug and the file.zip are the two parameter
func (s Store) UploadWatersheds(c echo.Context) error {
	wid, err := uuid.Parse(c.Param("watershed_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	r, err := models.UploadWatersheds(s.Connection, wid, file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, r)
}
