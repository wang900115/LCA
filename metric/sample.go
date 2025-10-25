package metric

import "slices"

// Hazen Method Interpolation to calculate percentiles.
func CalculatePercentiles(values []int64, ps []float64) []float64 {
	scores := make([]float64, len(ps))
	size := len(values)
	if size == 0 {
		return scores
	}
	slices.Sort(values)
	for i, p := range ps {
		pos := p * float64(size+1)
		if pos < 1.0 {
			scores[i] = float64(values[0])
		} else if pos >= float64(size) {
			scores[i] = float64(values[size-1])
		} else {
			lower := values[int(pos)-1]
			upper := values[int(pos)]
			fraction := pos - float64(int(pos))
			scores[i] = float64(lower) + fraction*float64(upper-lower)
		}
	}
	return scores
}
