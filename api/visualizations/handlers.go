package visualizations

import (
	"net/http"
	"strings"

	"github.com/USACE/water-api/api/messages"
	"github.com/USACE/water-api/api/visualizations/models"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/labstack/echo/v4"
)

func (s Store) ListVisualizations(c echo.Context) error {

	vv, err := models.ListVisualizations(s.Connection)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, vv)
}

func (s Store) GetVisualization(c echo.Context) error {

	vSlug := c.Param("visualization_slug")
	t, err := models.GetVisualization(s.Connection, &vSlug)

	if err != nil {
		if pgxscan.NotFound(err) {
			return c.JSON(http.StatusNotFound, messages.DefaultMessageNotFound)
		}
		// return c.JSON(http.StatusInternalServerError, messages.DefaultMessageInternalServerError)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, t)
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
