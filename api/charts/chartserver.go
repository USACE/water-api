package charts

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type (
	// ChartServer holds properties for a chart server
	// .ChartMap holds the same information as .Charts
	// .ChartMap is indexed by Slug string for fast lookups
	ChartServer struct {
		URL      *url.URL
		Charts   []ChartInfo
		ChartMap map[string]ChartInfo
	}

	ChartServerConfig struct {
		URLString string
	}

	ChartInfo struct {
		Slug           string   `json:"slug"`
		URL            string   `json:"url"`
		Name           string   `json:"name"`
		Description    string   `json:"description"`
		RequiredParams []string `json:"required_params"`
	}
)

func NewChartServer(cfg ChartServerConfig) (*ChartServer, error) {
	// Parse URL
	u, err := url.Parse(cfg.URLString)
	if err != nil {
		return nil, err
	}

	cs := ChartServer{
		URL:      u,
		Charts:   make([]ChartInfo, 0),
		ChartMap: make(map[string]ChartInfo),
	}

	// Fetch list of charts supported by chartserver hosted at url in config
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &cs.Charts); err != nil {
		return nil, err
	}

	// Convert list of charts to map indexed by slug for fast lookups
	for _, item := range cs.Charts {
		cs.ChartMap[item.Slug] = item
	}

	return &cs, nil
}
