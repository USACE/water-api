package chartserver

import (
	"net/url"
)

type Config struct {
	URLString string
}

type ChartServer struct {
	URL *url.URL
}

func NewChartServer(cfg Config) (*ChartServer, error) {
	u, err := url.Parse(cfg.URLString)
	if err != nil {
		return nil, err
	}
	return &ChartServer{URL: u}, nil
}
