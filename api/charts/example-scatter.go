package charts

import (
	"fmt"
	"net/url"
)

// ExampleScatter thinly wraps ChartDetail to allow implementing
// the Renderer interface
type ExampleScatter struct {
	ChartDetail *ChartDetail
}

func (d ExampleScatter) ChartTypeSlug() string {
	return d.ChartDetail.Type
}

func (d ExampleScatter) QueryValues() url.Values {
	q := url.Values{}
	for _, m := range d.ChartDetail.Mapping {
		if m.LatestValue != nil {
			latest := *m.LatestValue
			_, v := latest[0], latest[1]

			switch m.Variable {
			case "pointcount":
				q.Add(m.Variable, fmt.Sprintf("%v", v))
			}
		}
	}
	return q
}
