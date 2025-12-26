package services

import "math"

/*
Package services contains domain logic.

metrics.go implements deterministic computations over integer sequences.
*/

// ComputeMetrics calculates metrics for the given sequence.
// - SumFourthPowersNonPos: sum(yn^4) for yn <= 0
// - Min/Max/Count
func ComputeMetrics(values []int64) (count int, sumFourthNonPos int64, min int64, max int64) {
	count = len(values)
	if count == 0 {
		return 0, 0, 0, 0
	}

	min = values[0]
	max = values[0]

	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
		if v <= 0 {
			// v^4 can overflow int64 for large v; in this project we keep it simple.
			// If you expect large magnitudes, add overflow checks or use big.Int.
			p := int64(math.Pow(float64(v), 4))
			sumFourthNonPos += p
		}
	}
	return
}
