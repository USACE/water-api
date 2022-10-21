package ratings

import (
	"net/http"
	"strconv"

	"github.com/USACE/water-api/api/ratings/models"
	"github.com/labstack/echo/v4"
)

// Rate1Var
func (s Store) Rate1Var(c echo.Context) error {
	slug := c.Param("slug")
	method:=c.Param("method")
	ival, err := strconv.ParseFloat(c.Param("ival"), 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	dep, err := models.Rate1Var(s.Connection, slug, method, ival)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, dep)
}
