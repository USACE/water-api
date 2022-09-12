package chartserver

import (
	"github.com/USACE/water-api/api/helpers"
)

type DamProfileChartInput struct {
	Pool    float64 `querystring:"pool"`
	Tail    float64 `querystring:"tail"`
	Inflow  float64 `querystring:"inflow"`
	Outflow float64 `querystring:"outflow"`
}

func (s ChartServer) DamProfileChart(input DamProfileChartInput) (string, error) {
	u := *s.URL
	u.Path = u.Path + "/example-scatter"                     // Build URL Path
	u.RawQuery = helpers.StructToQueryValues(input).Encode() // Build URL Query Params
	return s.Get(&u)
}
