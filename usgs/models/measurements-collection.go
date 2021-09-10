package models

import (
	"encoding/json"
	"errors"

	"github.com/USACE/water-api/helpers"
)

// MeasurementCollection is an array of Measurements
type MeasurementCollection struct {
	Items []Measurement
}

// MeasurementCollection.UnmarshalJSON is a composition of the MeasurementCollection struct
func (c *MeasurementCollection) UnmarshalJSON(b []byte) error {
	switch helpers.JSONType(b) {
	case "ARRAY":
		return json.Unmarshal(b, &c.Items)
	case "OBJECT":
		c.Items = make([]Measurement, 1)
		return json.Unmarshal(b, &c.Items[0])
	default:
		return errors.New("payload not recognized as JSON array or object")
	}
}

// ParameterMeasurementCollection is an array of ParameterMeasurements
type ParameterMeasurementCollection struct {
	SiteNumber string `param:"site_number"`
	Items      []ParameterMeasurements
}

// ParameterMeasurementCollection.UnmarshalJSON is a compostion of the ParameterMeasurementCollection struct
func (c *ParameterMeasurementCollection) UnmarshalJSON(b []byte) error {
	switch helpers.JSONType(b) {
	case "ARRAY":
		return json.Unmarshal(b, &c.Items)
	case "OBJECT":
		c.Items = make([]ParameterMeasurements, 1)
		return json.Unmarshal(b, &c.Items[0])
	default:
		return errors.New("payload not recognized as JSON array or object")
	}
}
