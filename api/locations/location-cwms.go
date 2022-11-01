package locations

import (
	"fmt"
)

type (

	// CwmsLocation holds CwmsLocation-specific attributes and a pointer to the underlying Location
	// It allows for specific implementation of CRUD Interfaces
	CwmsLocation struct {
		Info                   LocationInfo // private attribute
		CwmsLocationAttributes CwmsLocationAttributes
	}

	// CwmsLocation holds CwmsLocation-specific attributes
	CwmsLocationAttributes struct {
		PublicName string
		Kind       string
	}
)

func NewCwmsLocation(l LocationInfo) (CwmsLocation, error) {

	cla, err := NewCwmsLocationAttributes(l.Attributes)
	if err != nil {
		return CwmsLocation{}, fmt.Errorf("[%v]: error creating CwmsLocationAttributes from provided information; %s", l, err.Error())
	}

	return CwmsLocation{Info: l, CwmsLocationAttributes: cla}, nil
}

func NewCwmsLocationAttributes(la LocationAttributes) (CwmsLocationAttributes, error) {

	var cla CwmsLocationAttributes

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
// Methods to Support Location Interface
// ////////////////////////////////////////////

func (cl CwmsLocation) LocationInfo() *LocationInfo {
	li := LocationInfo(cl.Info)
	li.Attributes = LocationAttributes{
		"public_name": cl.CwmsLocationAttributes.PublicName,
		"kind":        cl.CwmsLocationAttributes.Kind,
	}
	return &li
}
