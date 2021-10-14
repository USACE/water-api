package cwms

import (
	"net/http"

	"github.com/USACE/water-api/api/cwms/models"
	"github.com/USACE/water-api/api/helpers"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ListLevelKind
func (s Store) ListLevelKind(c echo.Context) error {
	lk, err := models.ListLevelKind(s.Connection)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, lk)
}

// CreateLevelKind
func (s Store) CreateLevelKind(c echo.Context) error {
	name := c.Param("name")
	// Get unique slug
	slug_name, err := helpers.NextUniqueSlug(s.Connection, "level_kind", "slug", name, "", "")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	lvl, err := models.CreateLevelKind(s.Connection, slug_name, name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, lvl)
}

// DeleteLevelKind
func (s Store) DeleteLevelKind(c echo.Context) error {
	ls := c.Param("slug")
	count, err := models.DeleteLevelKind(s.Connection, ls)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	} else if count == 0 {
		return c.JSON(http.StatusNotFound, "No record found for "+ls)
	}
	return c.JSON(http.StatusOK, "DELETED: "+ls)
}

// CreateLocationLevels
func (s Store) CreateLocationLevels(c echo.Context) error {
	var p []models.Levels
	if err := (&echo.DefaultBinder{}).BindBody(c, &p); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"Fail": err.Error()})
	}
	loc_id, _ := uuid.Parse(c.Param("location_id"))
	err := models.CreateLocationLevels(s.Connection, loc_id, p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"Fail": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"Success": "Values saved"})
}

// UpdateLocationLevels
func (s Store) UpdateLocationLevels(c echo.Context) error {
	var p []models.Levels
	if err := (&echo.DefaultBinder{}).BindBody(c, &p); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"Fail": err.Error()})
	}
	loc_id, _ := uuid.Parse(c.Param("location_id"))
	err := models.UpdateLocationLevels(s.Connection, loc_id, p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"Fail": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"Success": "Values updated"})
}

// ListLevelValues
func (s Store) ListLevelValues(c echo.Context) error {
	ls := c.Param("location_slug")
	lk := c.Param("level_kind")
	lvls, err := models.ListLevelValues(s.Connection, ls, lk)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	var rtn = struct {
		LocationSlug string               `json:"location_slug"`
		LevelKind    string               `json:"level_kind"`
		Levels       []models.LevelValues `json:"levels"`
	}{LocationSlug: ls, LevelKind: lk, Levels: lvls}

	return c.JSON(http.StatusOK, rtn)
}
