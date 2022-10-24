package models

type CwmsLocation struct {
	Location          // abstract location fields
	Name       string `json:"name"`
	PublicName string `json:"public_name"`
	Kind       string `json:"kind"`
}
