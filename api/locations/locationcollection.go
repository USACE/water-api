package locations

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/USACE/water-api/api/helpers"
)

type (
	// LocationCollection supports marshal/unmarshal of Locations passed either a single
	// struct (Object) or array of objects (slice)
	LocationCollection struct {
		Items []Location `json:"items"`
	}
)

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

func (lc LocationCollection) LocationCreatorCollection() (LocationCreatorCollection, error) {

	cc := make([]LocationCreator, len(lc.Items))

	for idx, item := range lc.Items {
		switch item.DatasourceType {
		case "cwms-location":
			if l, err := NewCwmsLocation(&item); err != nil {
				return LocationCreatorCollection{}, err
			} else {
				cc[idx] = l
			}
		default:
			return LocationCreatorCollection{}, fmt.Errorf("CREATE not implemented for datasource_type=%s", item.DatasourceType)
		}

	}
	return LocationCreatorCollection{Items: cc}, nil
}
