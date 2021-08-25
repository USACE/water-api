package models

import (
	"encoding/json"
	"errors"

	"github.com/USACE/water-api/helpers"
)

type SiteParameterCollection struct {
	Items []SiteParameter `json:"items"`
}

func (c *SiteParameterCollection) UnmarshalJSON(b []byte) error {
	switch helpers.JSONType(b) {
	case "ARRAY":
		return json.Unmarshal(b, &c.Items)
	case "OBJECT":
		c.Items = make([]SiteParameter, 1)
		return json.Unmarshal(b, &c.Items[0])
	default:
		return errors.New("payload not recognized as JSON array or object")
	}
}
