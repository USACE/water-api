package chartserver

import (
	"io"
	"net/http"
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

func (s ChartServer) Get(u *url.URL) (string, error) {
	urlstr := u.String()
	resp, err := http.Get(urlstr)
	if err != nil {
		return "", err
	}
	// TODO; Check Response Status (200), Handle Non-200 status gracefully
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// Convert to string
	b := string(body)
	return b, nil
}
