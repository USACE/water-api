package models

import (
	"encoding/json"
	"errors"

	"github.com/USACE/water-api/api/helpers"
)

type Location struct {
	ProviderSlug string           `json:"provider_slug"`
	ProviderName string           `json:"provider_name"`
	Slug         string           `json:"slug"`
	Geometry     helpers.Geometry `json:"geometry"`
	State        *string          `json:"state"` // state abbreviation (e.g. MN, TN, WV, FL)
}

type LocationCollection struct {
	Items []Location `json:"items"`
}

func (c *LocationCollection) UnmarshalJSON(b []byte) error {
	switch helpers.JSONType(b) {
	case "ARRAY":
		return json.Unmarshal(b, &c.Items)
	case "OBJECT":
		c.Items = make([]Location, 1)
		return json.Unmarshal(b, &c.Items[0])
	default:
		return errors.New("payload not recognized as JSON array or object")
	}
}
