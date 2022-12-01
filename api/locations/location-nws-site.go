package locations

import (
	"fmt"
)

type (
	// NwsSite holds NwsSite-specific attributes and a pointer to the underlying Location
	// It allows for specific implementation of CRUD Interfaces
	NwsSite struct {
		Info              LocationInfo
		NwsSiteAttributes NwsSiteAttributes
	}

	// NwsSite holds NwsSite-specific attributes
	NwsSiteAttributes struct {
		Name string
	}
)

func NewNwsSite(l LocationInfo) (NwsSite, error) {

	a, err := NewNwsSiteAttributes(l.Attributes)
	if err != nil {
		return NwsSite{}, fmt.Errorf("[%v]: error creating NwsSiteAttributes from provided information; %s", l, err.Error())
	}

	return NwsSite{Info: l, NwsSiteAttributes: a}, nil
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

	return a, nil

}

// ////////////////////////////////////////////
// Methods to Support LocationCreator Interface
// ////////////////////////////////////////////
func (n NwsSite) LocationInfo() *LocationInfo {
	li := LocationInfo(n.Info)
	li.Attributes = LocationAttributes{
		"name": n.NwsSiteAttributes.Name,
	}
	return &li
}
