package chartserver

import (
	"fmt"

	"github.com/USACE/water-api/api/helpers"
)

type DamProfileChartInput struct {
	Pool    float64 `querystring:"pool"`
	Tail    float64 `querystring:"tail"`
	Inflow  float64 `querystring:"inflow"`
	Outflow float64 `querystring:"outflow"`
}

func (s *ChartServer) GetDamProfileChart(input DamProfileChartInput) (string, error) {

	URL1 := *s.URL
	pURL1 := &URL1
	pURL1.Path = pURL1.Path + "/example-scatter"

	qv := helpers.StructToQueryValues(input)
	for v := range qv {
		fmt.Println(v)
	}

	return s.Get(pURL1)
}
