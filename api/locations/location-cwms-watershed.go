package locations

import (
	"fmt"
)

type (

	// CwmsWatershed holds CwmsWatershed-specific attributes and a pointer to the underlying Location
	// It allows for specific implementation of CRUD Interfaces
	CwmsWatershed struct {
		Info                    LocationInfo // private attribute
		CwmsWatershedAttributes CwmsWatershedAttributes
	}

	// CwmsWatershedAttributes holds CwmsWatershed-specific attributes
	CwmsWatershedAttributes struct{}
)

func NewCwmsWatershed(l LocationInfo) (CwmsWatershed, error) {

	a, err := NewCwmsWatershedAttributes(l.Attributes)
	if err != nil {
		return CwmsWatershed{}, fmt.Errorf("[%v]: error creating CwmsWatershedAttributes from provided information; %s", l, err.Error())
	}

	return CwmsWatershed{Info: l, CwmsWatershedAttributes: a}, nil
}

func NewCwmsWatershedAttributes(la LocationAttributes) (CwmsWatershedAttributes, error) {

	return CwmsWatershedAttributes{}, nil

}

// ////////////////////////////////////////////
// Methods to Support Location Interface
// ////////////////////////////////////////////

func (cw CwmsWatershed) LocationInfo() *LocationInfo {
	li := LocationInfo(cw.Info)
	li.Attributes = LocationAttributes{}
	return &li
}
