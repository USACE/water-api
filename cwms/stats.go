package cwms

import (
	"net/http"

	"github.com/USACE/water-api/cwms/models"
	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
)

func (s Store) ListStatsStates(c echo.Context) error {
	ss, err := models.ListStatsStates(s.Connection)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ss)
}

func (s Store) GetStatsState(c echo.Context) error {
	identifier := c.Param("state_id")
	ss, err := models.GetStatsState(s.Connection, &identifier)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ss)
}

func (s Store) ListStatsOffices(c echo.Context) error {
	so, err := models.ListStatsOffices(s.Connection)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, so)
}

func (s Store) GetStatsOffice(c echo.Context) error {
	officeID, err := uuid.Parse(c.Param("office_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	so, err := models.GetStatsOffice(s.Connection, &officeID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, so)
}
