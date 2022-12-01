package charts

import (
	"encoding/json"
	"fmt"
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
		Charts   []ChartType
		ChartMap map[string]ChartType
	}

	ChartServerConfig struct {
		URLString string
	}

	ChartType struct {
		Slug           string   `json:"slug"`
		URL            string   `json:"url"`
		Name           string   `json:"name"`
		Description    string   `json:"description"`
		RequiredParams []string `json:"required_params"`
	}

	Renderer interface {
		ChartTypeSlug() string
		QueryValues() url.Values
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
		Charts:   make([]ChartType, 0),
		ChartMap: make(map[string]ChartType),
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

func (cs ChartServer) Render(r Renderer) (*string, error) {

	slug := r.ChartTypeSlug()
	ChartType, ok := cs.ChartMap[slug]
	if !ok {
		return nil, fmt.Errorf("chart type '%s' not available from chartserver: %s", slug, cs.URL.String())
	}

	// Build URL
	u, err := url.Parse(ChartType.URL)
	if err != nil {
		return nil, err
	}
	u.RawQuery = r.QueryValues().Encode()

	urlstr := u.String()
	resp, err := http.Get(urlstr)
	if err != nil {
		return nil, err
	}
	// TODO; Check Response Status (200), Handle Non-200 status gracefully
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	b := string(body)
	return &b, nil

}
