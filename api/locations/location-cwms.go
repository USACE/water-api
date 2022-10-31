package locations

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type (

	// CwmsLocation holds CwmsLocation-specific attributes and a pointer to the underlying Location
	// It allows for specific implementation of CRUD Interfaces
	CwmsLocation struct {
		Location               *Location
		CwmsLocationAttributes CwmsLocationAttributes
	}

	// CwmsLocation holds CwmsLocation-specific attributes
	CwmsLocationAttributes struct {
		Name       string
		PublicName string
		Kind       string
	}
)

func NewCwmsLocation(l *Location) (CwmsLocation, error) {

	cla, err := NewCwmsLocationAttributes(l.Attributes)
	if err != nil {
		return CwmsLocation{}, fmt.Errorf("[%v]: error creating CwmsLocationAttributes from provided information; %s", *l, err.Error())
	}

	return CwmsLocation{Location: l, CwmsLocationAttributes: cla}, nil
}

func NewCwmsLocationAttributes(la LocationAttributes) (CwmsLocationAttributes, error) {

	var cla CwmsLocationAttributes

	// Name Attribute
	a1, ok := la["name"]
	if !ok {
		return CwmsLocationAttributes{}, fmt.Errorf("attribute 'name' is required")
	}
	name, ok := a1.(string)
	if !ok {
		return CwmsLocationAttributes{}, fmt.Errorf("attribute 'name' must be a string")
	}
	cla.Name = name

	// Public Name
	a2, ok := la["public_name"]
	if !ok {
		return CwmsLocationAttributes{}, fmt.Errorf("attribute 'public_name' is required")
	}
	publicName, ok := a2.(string)
	if !ok {
		return CwmsLocationAttributes{}, fmt.Errorf("attribute 'public_name' must be a string")
	}
	cla.PublicName = publicName

	// Kind
	a3, ok := la["kind"]
	if !ok {
		return CwmsLocationAttributes{}, fmt.Errorf("attribute 'kind' is required")
	}
	kind, ok := a3.(string)
	if !ok {
		return CwmsLocationAttributes{}, fmt.Errorf("attribute 'kind' must be a string")
	}
	cla.Kind = kind

	return cla, nil

}

// ////////////////////////////////////////////
// Methods to Support LocationCreator Interface
// ////////////////////////////////////////////

func (l CwmsLocation) LocationInfo() *Location {
	return l.Location
}

func (l CwmsLocation) CreateAttributes(tx *pgx.Tx, locationID *uuid.UUID) error {
	t := *tx
	if _, err := t.Exec(
		context.Background(),
		`INSERT INTO cwms_location (location_id, name, public_name, kind_id) VALUES
		    ($1, $2, $3, (SELECT id FROM cwms_location_kind WHERE UPPER(name) = UPPER($4)))`,
		locationID, l.CwmsLocationAttributes.Name, l.CwmsLocationAttributes.PublicName, l.CwmsLocationAttributes.Kind,
	); err != nil {
		return err
	}
	return nil
}
