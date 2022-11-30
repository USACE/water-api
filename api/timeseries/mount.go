package timeseries

import (
	"github.com/USACE/water-api/api/app"
	"github.com/USACE/water-api/api/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

type Store struct {
	Connection *pgxpool.Pool
}

func Mount(conn *pgxpool.Pool, e *echo.Echo, config *app.Config) {

	s := Store{
		Connection: conn,
	}

	// Public Routes
	public := e.Group("")

	// Key Only Group
	key := e.Group("")
	key.Use(middleware.KeyAuth(config.ApplicationKey))

	// Timeseries
	public.GET("/timeseries", s.ListTimeseries)                       // LIST
	key.POST("/providers/:provider/timeseries", s.CreateTimeseries)   // CREATE
	key.PUT("/providers/:provider/timeseries", s.UpdateTimeseries)    // UPDATE
	key.DELETE("/providers/:provider/timeseries", s.DeleteTimeseries) // DELETE

	// Timeseries Values
	key.POST("/providers/:provider/timeseries/values", s.CreateOrUpdateTimeseriesValues) // VALUES

	// Timeseries Groups
	public.GET("/providers/:provider/timeseries_groups", s.ListTimeseriesGroups)                       // LIST GROUPS
	public.GET("/providers/:provider/timeseries_groups/:timeseries_group", s.GetTimeseriesGroupDetail) // GET GROUP
	// public.GET("/providers/:provider/timeseries_groups/:timeseries_group/values", s.GetTimeseriesGroupValues)      // GET VALUES (EXTRACT)
	key.POST("/providers/:provider/timeseries_groups", s.CreateTimeseriesGroups) // CREATE GROUP(S)
	// key.PUT("/providers/:provider/timeseries_groups/:timeseries_group", s.UpdateTimeseriesGroup)                   // UPDATE GROUP
	key.DELETE("/providers/:provider/timeseries_groups/:timeseries_group", s.DeleteTimeseriesGroup) // DELETE GROUP
	// key.POST("/providers/:provider/timeseries_groups/:timeseries_group/members", s.AddTimeseriesGroupMembers)      // ADD GROUP MEMBER(S)
	// key.DELETE("/providers/:provider/timeseries_groups/:timeseries_group/members", s.RemoveTimeseriesGroupMembers) // REMOVE GROUP MEMBER(S)

}
