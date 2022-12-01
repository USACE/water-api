package charts

import (
	"fmt"
	"net/url"
)

// DamProfileChartInput thinly wraps ChartDetail to allow implementing
// the Renderer interface
type DamProfileChart struct {
	ChartDetail *ChartDetail
}

func (d DamProfileChart) ChartTypeSlug() string {
	return d.ChartDetail.Type
}

func (d DamProfileChart) QueryValues() url.Values {
	q := url.Values{}
	for _, m := range d.ChartDetail.Mapping {
		if m.LatestValue != nil {
			latest := *m.LatestValue
			_, v := latest[0], latest[1]

			switch m.Variable {
			case "pool", "tail", "inflow", "outflow", "damtop", "dambottom":
				q.Add(m.Variable, fmt.Sprintf("%v", v))
			case "top-of-dam", "top-of-flood":
				q.Add("level", fmt.Sprintf("%s,%v", m.Variable, v))
			}
		}
	}
	return q
}
