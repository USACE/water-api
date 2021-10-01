package helpers

import (
	"fmt"
)

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

func (g Geometry) EWKT(precision int) string {
	return fmt.Sprintf("SRID=4326;POINT(%.[3]*[1]f %.[3]*[2]f)", g.Coordinates[0], g.Coordinates[1], precision)
}
