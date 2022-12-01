package locations

import (
	"fmt"
)

type (

	// UsgsSite holds UsgsSite-specific attributes and a pointer to the underlying Location
	// It allows for specific implementation of CRUD Interfaces
	UsgsSite struct {
		Info               LocationInfo
		UsgsSiteAttributes UsgsSiteAttributes
	}

	// UsgsSiteAttributes holds UsgsSite-specific attributes
	UsgsSiteAttributes struct {
		StationName string
		SiteType    string
	}
)

func NewUsgsSite(l LocationInfo) (UsgsSite, error) {

	a, err := NewUsgsSiteAttributes(l.Attributes)
	if err != nil {
		return UsgsSite{}, fmt.Errorf("[%v]: error creating UsgsSiteAttributes from provided information; %s", l, err.Error())
	}

	return UsgsSite{Info: l, UsgsSiteAttributes: a}, nil
}

func NewUsgsSiteAttributes(la LocationAttributes) (UsgsSiteAttributes, error) {

	var a UsgsSiteAttributes

	// station_name attribute
	a1, ok := la["station_name"]
	if !ok {
		return UsgsSiteAttributes{}, fmt.Errorf("attribute 'station_name' is required")
	}
	stationName, ok := a1.(string)
	if !ok {
		return UsgsSiteAttributes{}, fmt.Errorf("attribute 'station_name' must be a string")
	}
	a.StationName = stationName

	// site_type attribute
	a3, ok := la["site_type"]
	if !ok {
		return UsgsSiteAttributes{}, fmt.Errorf("attribute 'site_type' is required")
	}
	siteType, ok := a3.(string)
	if !ok {
		return UsgsSiteAttributes{}, fmt.Errorf("attribute 'site_type' must be a string")
	}
	a.SiteType = siteType

	return a, nil

}

// ////////////////////////////////////////////
// Methods to Support LocationCreator Interface
// ////////////////////////////////////////////

func (u UsgsSite) LocationInfo() *LocationInfo {
	li := LocationInfo(u.Info)
	li.Attributes = LocationAttributes{
		"type":         u.UsgsSiteAttributes.SiteType,
		"station_name": u.UsgsSiteAttributes.StationName,
	}
	return &li
}
