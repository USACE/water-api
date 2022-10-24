package models

type NwsSite struct {
	Location        // abstract location fields
	Name     string `json:"usgs_site_number"`
}
