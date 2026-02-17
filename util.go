package aprs

import (
	"fmt"
	"math"
)

// unexported utility functions

// z3p returns a 3-zeros-padded string representation of the given int
func z3p(i int) string {
	return fmt.Sprintf("%03d", i)
}

// decToDMS takes a float latitude or longitude and converts it to
// degrees, minutes (second as decimal), and a hemisphere string.
func decToDMS(l float64, hems [2]string) (float64, float64, string) {
	deg, frac := math.Modf(math.Abs(l))
	min := frac * 60.0

	if l < 0 {
		return deg, min, hems[1]
	}
	return deg, min, hems[0]
}
