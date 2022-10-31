package locations

import (
	"github.com/USACE/water-api/api/helpers"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type (
	LocationAttributes map[string]interface{}

	// Location contains standardized information for all locations,
	// regardless of the kind of location and is used to serialize/deserialize
	// location information from structs to JSON for all kinds of locations
	Location struct {
		Provider           string             `json:"provider"`
		ProviderName       string             `json:"provider_name"`
		DatasourceType     string             `json:"datasource_type"`
		DatasourceTypeName string             `json:"datasource_type_name"`
		Slug               string             `json:"slug"`
		Geometry           helpers.Geometry   `json:"geometry"`
		State              *string            `json:"state"`
		Attributes         LocationAttributes `json:"attributes"` // Non-Standard Attributes
	}
)

// ////////////////////////////////////////////
// Methods to Support LocationCreator Interface
// ////////////////////////////////////////////

// NOTE: Saving a Location saves parent record in the location table, but does not save additional attributes.
//
//	To save additional attributes, create a new type and implement the LocationCreator interface.
//	See CwmsLocation, NwsSite, UsgsSite structs as examples
//

func (l Location) LocationInfo() *Location {
	return &l
}

func (l Location) CreateAttributes(tx *pgx.Tx, locationID *uuid.UUID) error {
	return nil
}
