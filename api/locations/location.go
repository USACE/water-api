package locations

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/USACE/water-api/api/helpers"
)

type (
	// Location defines the behaviors that must be supported for all
	// location kinds to be treated like a location
	Location interface {
		LocationInfo() *LocationInfo
	}

	// LocationCreatorCollection holds LocationCreator interfaces
	// to support different behaviors for locations that have different
	// underlying datasource_type properties
	LocationCollection struct {
		Items []Location
	}

	LocationAttributes map[string]interface{}

	// LocationInfo contains standardized information for all locations,
	// regardless of the kind of location and is used to serialize/deserialize
	// location information from structs to JSON for all kinds of locations
	LocationInfo struct {
		Provider           string             `json:"provider"`
		ProviderName       string             `json:"provider_name"`
		DatasourceType     string             `json:"datasource_type"`
		DatasourceTypeName string             `json:"datasource_type_name"`
		Code               string             `json:"code"` // unique string; e.g. cwms-location: "name", usgs-site: "station number", nws-site: "nws_li"
		Slug               string             `json:"slug"`
		Geometry           helpers.Geometry   `json:"geometry"`
		State              *string            `json:"state"`
		Attributes         LocationAttributes `json:"attributes"` // Non-Standard Attributes
	}

	// LocationInfoCollection supports marshal/unmarshal of Locations passed either a single
	// struct (Object) or array of objects (slice)
	LocationInfoCollection struct {
		Items []LocationInfo `json:"items"`
	}
)

// UnmarshalJSON supports posting LocationInfo as an object or an array of objects
func (c *LocationInfoCollection) UnmarshalJSON(b []byte) error {
	switch helpers.JSONType(b) {
	case "ARRAY":
		return json.Unmarshal(b, &c.Items)
	case "OBJECT":
		c.Items = make([]LocationInfo, 1)
		return json.Unmarshal(b, &c.Items[0])
	default:
		return errors.New("payload not recognized as JSON array or object")
	}
}

func (lc LocationInfoCollection) LocationCollection() (LocationCollection, error) {

	empty := LocationCollection{
		Items: make([]Location, 0),
	}

	cc := make([]Location, len(lc.Items))
	for idx, item := range lc.Items {
		switch item.DatasourceType {
		case "cwms-location":
			if l, err := NewCwmsLocation(item); err != nil {
				return empty, err
			} else {
				cc[idx] = l
			}
		case "usgs-site":
			if l, err := NewUsgsSite(item); err != nil {
				return empty, err
			} else {
				cc[idx] = l
			}
		case "nws-site":
			if l, err := NewNwsSite(item); err != nil {
				return empty, err
			} else {
				cc[idx] = l
			}
		default:
			return LocationCollection{}, fmt.Errorf("CREATE not implemented for datasource_type=%s", item.DatasourceType)
		}
	}
	return LocationCollection{Items: cc}, nil
}
