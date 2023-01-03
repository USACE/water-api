package charts

import (
	"fmt"
	"net/url"
)

// DamProfileChart thinly wraps ChartDetail to allow implementing
// the Renderer interface
type DamProfileChart struct {
	ChartDetail *ChartDetail
}

func (d DamProfileChart) ChartTypeSlug() string {
	return d.ChartDetail.Type
}

func (d DamProfileChart) QueryValues() url.Values {

	LEVELNAMES := map[string]string{
		"level-top-of-normal":    "Top of Normal",
		"level-bottom-of-normal": "Bottom of Normal",
		"level-top-of-flood":     "Top of Flood",
	}

	q := url.Values{}
	for _, m := range d.ChartDetail.Mapping {
		if m.LatestValue != nil {
			latest := *m.LatestValue
			_, v := latest[0], latest[1]

			switch m.Variable {
			case "pool", "tail", "inflow", "outflow":
				q.Add(m.Variable, fmt.Sprintf("%v", v))
			case "level-top-of-normal", "level-bottom-of-normal", "level-top-of-flood":
				q.Add("level", fmt.Sprintf("%s,%v", LEVELNAMES[m.Variable], v))
			case "damtop", "dambottom":
				q.Add(m.Variable, fmt.Sprintf("%v", v))             // add required variables damtop, dambottom
				q.Add("level", fmt.Sprintf("%s,%v", m.Variable, v)) // add level markerlines
			}
		}
	}
	return q
}
