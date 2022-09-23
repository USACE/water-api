package cwms

import (
	"net/http"

	"github.com/USACE/water-api/api/chartserver"
	water "github.com/USACE/water-api/api/water/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (s Store) GetProfileChart(c echo.Context) error {
	// TODO; Get location information needed to call chartserver
	//       mock in the data at this time.
	//
	locationSlug := c.Param("location_slug")
	visualizationTypeId, _ := uuid.Parse("53da77d0-6550-4f02-abf8-4bcd1a596a7c")

	v, err := water.GetVisualizationByLocation(s.Connection, &locationSlug, &visualizationTypeId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	//var m map[string]float64

	// Create map to hold variable maps
	m := make(map[string]map[string]interface{})

	for _, entry := range *v.Mapping {
		// fmt.Printf("%#v\n", entry)

		// Each variable will have a map to hold items like:
		// key, latest_time, latest_value

		variable_map := make(map[string]interface{})

		variable_map["key"] = entry.Key

		if entry.LatestValue != nil {
			variable_map["latest_value"] = *entry.LatestValue
		}
		if entry.LatestTime != nil {
			variable_map["latest_time"] = *entry.LatestTime
		}

		// assign the variable map to the larger map of variables
		m[entry.Variable] = variable_map

		// record := reflect.Indirect(reflect.ValueOf(entry))
		// fieldName := record.Type().Field(i).Name
		// //fmt.Println("KV Pair: ", value.Type().Field(0).Name, v)
		// // fmt.Println(fieldName, "->", entry)

	}

	// for v, x := range m {
	// 	fmt.Println(v, " -> ", x)
	// }

	pool_val := m["pool"]["latest_value"]
	if pool_val == nil {
		pool_val = float64(0)
	}

	inflow_val := m["inflow"]["latest_value"]
	if inflow_val == nil {
		inflow_val = float64(0)
	}

	outflow_val := m["outflow"]["latest_value"]
	if outflow_val == nil {
		outflow_val = float64(0)
	}

	damTop_val := m["top-of-dam"]["latest_value"]
	if damTop_val == nil {
		damTop_val = float64(0)
	}

	input := chartserver.DamProfileChartInput{
		Pool:      pool_val.(float64),
		Tail:      12.0,
		Inflow:    inflow_val.(float64),
		Outflow:   outflow_val.(float64),
		DamTop:    damTop_val.(float64),
		DamBottom: 300,
	}
	chart, err := s.ChartServer.DamProfileChart(input)
	if err != nil {
		return c.String(500, err.Error())
	}

	return c.HTML(200, chart)
}
