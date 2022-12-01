package timeseries

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/USACE/water-api/api/messages"
	"github.com/USACE/water-api/api/timeseries/models"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/labstack/echo/v4"
)

/////////////
// TIMESERIES
/////////////

func (s Store) ListTimeseries(c echo.Context) error {
	// Get filters from query provider= or datatype=
	var f models.TimeseriesFilter
	if err := c.Bind(&f); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	ll, err := models.ListTimeseries(s.Connection, &f)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ll)
}

func (s Store) CreateTimeseries(c echo.Context) error {
	var tsc models.TimeseriesCollection
	if err := c.Bind(&tsc); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// ensure provider slug for all posted locations matches route param :provider
	// if not, return status unauthorized
	routeProvider := c.Param("provider")
	for _, item := range tsc.Items {
		if !strings.EqualFold(routeProvider, item.Provider) {
			return c.String(
				http.StatusBadRequest,
				"timeseries in post body has provider that does not match route param :provider",
			)
		}
	}

	tt, err := tsc.Create(s.Connection, strings.ToLower(routeProvider))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, messages.DefaultMessageInternalServerError)
	}

	// If 0 new timeseries were created, return a RESTful 200
	if len(tt) == 0 {
		return c.JSON(http.StatusOK, tt)
	}

	// If at least 1 timeseries was created, return 201 with array of new timeseries
	return c.JSON(http.StatusCreated, tt)
}

func (s Store) CreateOrUpdateTimeseriesValues(c echo.Context) error {
	var tsc models.TimeseriesCollection
	if err := c.Bind(&tsc); err != nil {
		if strings.Contains(err.Error(), "parsing time") {
			return c.JSON(http.StatusBadRequest, messages.NewMessage("incorrect time format, use YYYY-MM-DDTHH:MM:SS-HH:MM"))
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	tt, err := tsc.CreateOrUpdateTimeseriesValues(s.Connection)
	if err != nil {
		if strings.Contains(err.Error(), "no records updated") {
			return c.JSON(http.StatusBadRequest, messages.NewMessage("values not provided in proper format"))
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusAccepted, tt)

}

func (s Store) UpdateTimeseries(c echo.Context) error {
	var tsc models.TimeseriesCollection
	if err := c.Bind(&tsc); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// ensure provider slug for all posted locations matches route param :provider
	// if not, return status unauthorized
	routeProvider := c.Param("provider")
	for _, item := range tsc.Items {
		if !strings.EqualFold(routeProvider, item.Provider) {
			return c.String(
				http.StatusBadRequest,
				"timeseries in post body has provider that does not match route param :provider",
			)
		}
	}

	tt, err := tsc.Update(s.Connection)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, messages.DefaultMessageInternalServerError)
	}

	return c.JSON(http.StatusOK, tt)
}

func (s Store) DeleteTimeseries(c echo.Context) error {
	var tsc models.TimeseriesCollection
	if err := c.Bind(&tsc); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// ensure provider slug for all posted locations matches route param :provider
	// if not, return status unauthorized
	routeProvider := c.Param("provider")
	for _, item := range tsc.Items {
		if !strings.EqualFold(routeProvider, item.Provider) {
			return c.String(
				http.StatusBadRequest,
				"timeseries in post body has provider that does not match route param :provider",
			)
		}
	}

	if err := tsc.Delete(s.Connection); err != nil {
		return c.JSON(http.StatusInternalServerError, messages.DefaultMessageInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{})
}

////////////////////
// TIMESERIES GROUPS
////////////////////

func (s Store) ListTimeseriesGroups(c echo.Context) error {

	var f models.TimeseriesGroupFilter
	if err := c.Bind(&f); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	gg, err := models.ListTimeseriesGroups(s.Connection, &f)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, gg)
}

func (s Store) GetTimeseriesGroupDetail(c echo.Context) error {

	var f models.TimeseriesGroupFilter
	if err := c.Bind(&f); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	g, err := models.GetTimeseriesGroupDetail(s.Connection, &f)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, &g)
}

func (s Store) CreateTimeseriesGroups(c echo.Context) error {

	var gc models.TimeseriesGroupCollection
	if err := c.Bind(&gc); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	routeProvider := c.Param("provider")
	for _, item := range gc.Items {
		// 1. Verify provider in payload body matches route param :provider
		if !strings.EqualFold(routeProvider, item.Provider) {
			return c.JSON(
				http.StatusBadRequest,
				messages.NewMessage(
					fmt.Sprintf(
						"timeseries group in post body has provider (%s) that does not match route param :provider (%s)",
						item.Provider,
						routeProvider,
					),
				),
			)
		}
	}
	gg, err := models.CreateTimeseriesGroups(s.Connection, &gc)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, gg)
}

func (s Store) DeleteTimeseriesGroup(c echo.Context) error {
	provider, slug := c.Param("provider"), c.Param("timeseries_group")

	if err := models.DeleteTimeseriesGroup(s.Connection, &provider, &slug); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{})
}

func (s Store) AddTimeseriesGroupMembers(c echo.Context) error {
	provider, slug := c.Param("provider"), c.Param("timeseries_group")
	var mc models.TimeseriesGroupMemberCollection
	if err := c.Bind(&mc); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	g, err := models.AddTimeseriesGroupMembers(s.Connection, &provider, &slug, &mc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, g)
}

func (s Store) RemoveTimeseriesGroupMembers(c echo.Context) error {

	provider, slug := c.Param("provider"), c.Param("timeseries_group")
	var mc models.TimeseriesGroupMemberCollection
	if err := c.Bind(&mc); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	g, err := models.RemoveTimeseriesGroupMembers(s.Connection, &provider, &slug, &mc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, g)
}

func (s Store) GetTimeseriesGroupValues(c echo.Context) error {

	// Uniquely Identify Timeseries Group
	provider, slug := c.Param("provider"), c.Param("timeseries_group")

	// Time Window
	after, before := c.QueryParam("after"), c.QueryParam("before")
	if (after == "" || before == "") && (after != before) {
		return c.JSON(http.StatusBadRequest, messages.NewMessage("query parameters 'after' and 'before' must both be provided, or both must left blank"))
	}

	var tw models.TimeWindow
	if err := models.NewTimeWindow(&tw, after, before); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	rows, err := s.Connection.Query(
		context.Background(),
		`WITH values AS (
			SELECT t.id    AS timeseries_id,
				   v.time  AS time,
				   v.value AS value
			FROM v_timeseries_group g
			JOIN timeseries_group_members m ON m.timeseries_group_id = g.id
			JOIN timeseries               t ON t.id = m.timeseries_id 
			JOIN timeseries_value         v ON v.timeseries_id = m.timeseries_id
			 AND v.time >= $3
			 AND v.time <= $4
		    WHERE g.provider = LOWER($1) AND g.slug = LOWER($2)
		    ORDER BY t.id, v.time
	   ), values_agg AS (
		    SELECT timeseries_id,
			       json_agg(json_build_array(time, value)) AS values
		    FROM values
		    GROUP BY timeseries_id
		)
		SELECT t.provider,
		       t.datatype,
			   t.key,
	           v.values
	    FROM values_agg v
	    JOIN v_timeseries t ON t.id = v.timeseries_id`,
		provider, slug, tw.After, tw.Before,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	defer rows.Close()

	type T struct {
		Provider string           `json:"provider"`
		Datatype string           `json:"datatype"`
		Key      string           `json:"key"`
		Values   *[][]interface{} `json:"values"` // may be empty [] or [["2022-09-27T12:00:00-05:00", 888.00], ["2022-09-27T13:00:00-05:00", 888.15]]
	}

	enc := json.NewEncoder(c.Response())
	for rows.Next() {
		var t T
		if err := pgxscan.ScanRow(&t, rows); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		if err := enc.Encode(t); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		c.Response().Flush()
	}

	// Any errors encountered by rows.Next or rows.Scan will be returned here
	if rows.Err() != nil {
		return err
	}

	// No errors found
	return nil
}
