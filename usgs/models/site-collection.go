package models

import (
	"encoding/json"
	"errors"

	"github.com/USACE/water-api/helpers"
)

type SiteCollection struct {
	Items []Site `json:"items"`
}

func (c *SiteCollection) UnmarshalJSON(b []byte) error {
	switch helpers.JSONType(b) {
	case "ARRAY":
		return json.Unmarshal(b, &c.Items)
	case "OBJECT":
		c.Items = make([]Site, 1)
		return json.Unmarshal(b, &c.Items[0])
	default:
		return errors.New("payload not recognized as JSON array or object")
	}
}
