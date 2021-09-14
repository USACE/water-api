package nws

import (
	"net/http"

	"github.com/USACE/water-api/messages"
	"github.com/USACE/water-api/nws/models"
	"github.com/georgysavva/scany/pgxscan"

	"github.com/labstack/echo/v4"
)

// ListNwsStages returns an array of NWS Location Stages
func (s Store) ListNwsStages(c echo.Context) error {
	ns, err := models.ListNwsStages(s.Connection)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ns)

}

// GetNwsStages returns a single NWS Stages Record
func (s Store) GetNwsStages(c echo.Context) error {
	id := c.Param("nwsid")
	ns, err := models.GetNwsStages(s.Connection, &id)
	if err != nil {
		if pgxscan.NotFound(err) {
			return c.JSON(http.StatusNotFound, messages.DefaultMessageNotFound)
		}
		return c.JSON(http.StatusInternalServerError, messages.DefaultMessageInternalServerError)
	}
	return c.JSON(http.StatusOK, ns)
}

// CreateNwsStage creates a new NWS Stages Record
func (s Store) CreateNwsStage(c echo.Context) error {
	var ns models.NwsStages
	if err := c.Bind(&ns); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	newNwsStage, err := models.CreateNwsStage(s.Connection, &ns)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, newNwsStage)
}

// UpdateNwsStages updates a single NWS Stages Record
func (s Store) UpdateNwsStages(c echo.Context) error {
	// nwsid from route params
	nwsid := c.Param("nwsid")
	if nwsid == "" {
		return c.String(http.StatusBadRequest, "nwsid not provided in route url")
	}
	// Payload
	var ns models.NwsStages
	if err := c.Bind(&ns); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// Check route params v. payload
	if nwsid != ns.NwsID {
		return c.String(http.StatusBadRequest, "nwsid in URL does not match request body")
	}
	sUpdated, err := models.UpdateNwsStages(s.Connection, &ns)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, sUpdated)

}
