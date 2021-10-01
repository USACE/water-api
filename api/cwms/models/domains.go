package models

import "github.com/google/uuid"

type Domain struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
