package interpolation

import "errors"

type InterpAccel struct {
	Cache     uint64
	MissCount uint64
	HitCount  uint64
}

// InterpFind
func InterpSearch(x_array []float64, x float64, index_lo uint64, index_hi uint64) uint64 {
	ilo := index_lo
	ihi := index_hi
	for ihi > ilo+1 {
		i := (ihi + ilo) / 2
		if x_array[i] > x {
			ihi = i
		} else {
			ilo = i
		}
	}
	return ilo
}

// InterpFind
func InterpFind(a *InterpAccel, x_array []float64, x float64) uint64 {
	x_index := a.Cache

	if x < x_array[x_index] {
		a.MissCount++
		a.Cache = InterpSearch(x_array, x, 0, x_index)
	} else if x >= x_array[x_index+1] {
		a.MissCount++
		a.Cache = InterpSearch(x_array, x, x_index, uint64(len(x_array)-1))
	} else {
		a.HitCount++
	}

	return a.Cache
}

// Linear
func LinearEval(x_array []float64, y_array []float64, x float64, a *InterpAccel) (float64, error) {
	var (
		x_lo, y_lo, x_hi, y_hi, dx float64
		index                      uint64
	)

	if a == (&InterpAccel{}) {
		index = InterpFind(a, x_array, x)
	} else {
		index = InterpSearch(x_array, x, 0, uint64(len(x_array)-1))
	}
	x_lo = x_array[index]
	x_hi = x_array[index+1]
	y_lo = y_array[index]
	y_hi = y_array[index+1]
	dx = x_hi - x_lo
	if dx > 0.0 {
		return y_lo + (x-x_lo)/dx*(y_hi-y_lo), nil
	}

	return 0.0, errors.New("dx not greater than 0.0")
}
