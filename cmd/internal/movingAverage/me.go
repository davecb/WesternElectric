package movingAverage

/*
 * Moving Average - the thing to compare outliers against
 * Derived from Snowball @ stackoverflow's
 */

// NewFMovingAverage - return a new moving-average function
func NewFMovingAverage(nSamples int) func(float64) float64 {
	bins := make([]float64, nSamples)
	var average = 0.0
	i := 0
	return func(x float64) float64 {
		average += (x - bins[i]) / float64(nSamples)
		bins[i] = x
		i = (i + 1) % nSamples
		return average
	}
}
