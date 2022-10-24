package models

type UsgsSite struct {
	Location           // abstract location fields
	SiteNumber  string `json:"usgs_site_number"`
	StationName string `json:"station_name"`
	SiteType    string `json:"site_type"`
}
