package water

import (
	"net/http"
	"strings"

	"github.com/USACE/water-api/api/messages"
	"github.com/USACE/water-api/api/water/models"
	"github.com/labstack/echo/v4"
)

func (s Store) ListVisualizations(c echo.Context) error {

	vv, err := models.ListVisualizations(s.Connection)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, vv)
}

// CreateVisualization creates a single new Visualization
func (s Store) CreateVisualization(c echo.Context) error {

	var v models.Visualization
	if err := c.Bind(&v); err != nil {
		return c.JSON(http.StatusBadRequest, messages.DefaultMessageBadRequest)
	}
	vNew, err := models.CreateVisualization(s.Connection, &v)
	if err != nil {
		if strings.Contains(err.Error(), `null value in column "location_id"`) {
			return c.JSON(http.StatusBadRequest, messages.NewMessage("invalid location_slug"))
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, &vNew)

}

func (s Store) CreateOrUpdateVisualizationMapping(c echo.Context) error {

	slug := c.Param("visualization_slug")

	var vm models.VisualizationMappingCollection
	if err := c.Bind(&vm); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	_, err := models.CreateOrUpdateVisualizationMapping(s.Connection, vm, &slug)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusAccepted)

}
