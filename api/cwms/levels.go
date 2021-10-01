package cwms

import (
	"net/http"

	"github.com/USACE/water-api/api/cwms/models"
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

// CreateLocationLevel
// func (s Store) CreateLocationLevel(c echo.Context) error {
// 	var lvl models.Level
// 	lvl, err := models.CreateLocationLevels(s.Connection)
// 	return c.JSON(http.StatusOK, &lvl)
// }
