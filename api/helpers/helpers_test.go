package helpers

import (
	"testing"
)

func TestStructToQueryValues(t *testing.T) {

	type testStruct struct {
		Pool    float64 `querystring:"pool"`
		Tail    float64 `querystring:"tail"`
		Inflow  float64 `querystring:"inflow"`
		Outflow float64 `querystring:"outflow"`
	}

	s := testStruct{
		Pool:    15,
		Tail:    12.02,
		Inflow:  1200,
		Outflow: 600,
	}

	q := StructToQueryValues(s)
	if len(q) != 4 {
		t.Errorf("length of map is %d, should be 4", len(q))
	}
}
