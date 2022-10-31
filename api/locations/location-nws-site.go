package locations

import (
	"context"
	"fmt"

	"github.com/USACE/water-api/api/helpers"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type (
	// NwsSite holds NwsSite-specific attributes and a pointer to the underlying Location
	// It allows for specific implementation of CRUD Interfaces
	NwsSite struct {
		Location          Location
		NwsSiteAttributes NwsSiteAttributes
	}

	// NwsSite holds NwsSite-specific attributes
	NwsSiteAttributes struct {
		Name  string
		NwsLI string // nws location identifier or "handbook 5 code"
	}
)

func NewNwsSite(l Location) (NwsSite, error) {

	a, err := NewNwsSiteAttributes(l.Attributes)
	if err != nil {
		return NwsSite{}, fmt.Errorf("[%v]: error creating NwsSiteAttributes from provided information; %s", l, err.Error())
	}

	l.Slug = helpers.Slugify(a.NwsLI) // enforce nws_li used in slug creation

	return NwsSite{Location: l, NwsSiteAttributes: a}, nil
}

func NewNwsSiteAttributes(la LocationAttributes) (NwsSiteAttributes, error) {

	var a NwsSiteAttributes

	// Name Attribute
	a1, ok := la["name"]
	if !ok {
		return NwsSiteAttributes{}, fmt.Errorf("attribute 'name' is required")
	}
	name, ok := a1.(string)
	if !ok {
		return NwsSiteAttributes{}, fmt.Errorf("attribute 'name' must be a string")
	}
	a.Name = name

	// Public Name
	a2, ok := la["nws_li"]
	if !ok {
		return NwsSiteAttributes{}, fmt.Errorf("attribute 'nws_li' is required")
	}
	nwsLI, ok := a2.(string)
	if !ok {
		return NwsSiteAttributes{}, fmt.Errorf("attribute 'nws_li' must be a string")
	}
	a.NwsLI = nwsLI

	return a, nil

}

// ////////////////////////////////////////////
// Methods to Support LocationCreator Interface
// ////////////////////////////////////////////

func (l NwsSite) LocationInfo() Location {
	return l.Location
}

func (l NwsSite) CreateAttributes(tx *pgx.Tx, locationID *uuid.UUID) error {
	t := *tx
	if _, err := t.Exec(
		context.Background(),
		`INSERT INTO nws_site (location_id, name, nws_li) VALUES ($1, $2, $3)`,
		locationID, l.NwsSiteAttributes.Name, l.NwsSiteAttributes.NwsLI,
	); err != nil {
		return err
	}
	return nil
}
