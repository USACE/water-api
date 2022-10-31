package locations

import (
	"context"
	"fmt"

	"github.com/USACE/water-api/api/helpers"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type (

	// UsgsSite holds UsgsSite-specific attributes and a pointer to the underlying Location
	// It allows for specific implementation of CRUD Interfaces
	UsgsSite struct {
		Location           Location
		UsgsSiteAttributes UsgsSiteAttributes
	}

	// UsgsSiteAttributes holds UsgsSite-specific attributes
	UsgsSiteAttributes struct {
		StationName string
		SiteNumber  string
		SiteType    string
	}
)

func NewUsgsSite(l Location) (UsgsSite, error) {

	a, err := NewUsgsSiteAttributes(l.Attributes)
	if err != nil {
		return UsgsSite{}, fmt.Errorf("[%v]: error creating UsgsSiteAttributes from provided information; %s", l, err.Error())
	}

	l.Slug = helpers.Slugify(a.SiteNumber) // enforce nws_li attribute used in slug creation

	return UsgsSite{Location: l, UsgsSiteAttributes: a}, nil
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

	// site_number attribute
	a2, ok := la["site_number"]
	if !ok {
		return UsgsSiteAttributes{}, fmt.Errorf("attribute 'site_number' is required")
	}
	siteNumber, ok := a2.(string)
	if !ok {
		return UsgsSiteAttributes{}, fmt.Errorf("attribute 'site_number' must be a string")
	}
	a.SiteNumber = siteNumber

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

func (l UsgsSite) LocationInfo() Location {
	return l.Location
}

func (l UsgsSite) CreateAttributes(tx *pgx.Tx, locationID *uuid.UUID) error {
	t := *tx
	if _, err := t.Exec(
		context.Background(),
		`INSERT INTO usgs_site (location_id, station_name, site_number, site_type_id) VALUES
		    ($1, $2, $3, (SELECT id FROM usgs_site_type WHERE UPPER(abbreviation) = UPPER($4)))`,
		locationID, l.UsgsSiteAttributes.StationName, l.UsgsSiteAttributes.SiteNumber, l.UsgsSiteAttributes.SiteType,
	); err != nil {
		return err
	}
	return nil
}
