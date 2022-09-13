package movingAverage

import (
	"math"
)

/*
 * Moving Average -- this is a simple (non-exponential) moving average, using
 * the Knuth-Welford algorithm, as described in
 * https://pkg.go.dev/github.com/eclesh/welford#section-readme -- does an infinite sequence in go
 * https://stackoverflow.com/questions/1174984/how-to-efficiently-calculate-a-running-standard-deviation
 * https://www.johndcook.com/blog/2008/09/26/comparing-three-methods-of-computing-standard-deviation/
 * https://www.johndcook.com/blog/standard_deviation/
 *
 * The proximate code is from Tobia Vendrame's python example at
 * https://medium.com/@tobMTV/19-streaming-algorithm-for-the-mean-and-the-variance-b8a9d8faf8a6

   variance(samples):
		M := 0
		S := 0
		for k from 1 to N:
			x := samples[k]
			oldM := M
			M := M + (x-M)/k
			S := S + (x-M)*(x-oldM)
		return S/(N-1)

*
*/

// New takes a number of sample to consider and generates a
// function, conventionally called "add", that computes a moving average
// and standard deviation as each sample is added to the sample set.
func New(nSamples int) func(s float64) (float64, float64) {
	var i int

	bins := make([]float64, nSamples)
	return func(new float64) (float64, float64) {
		var Mean, S float64 // S is the accumulator for the variance and SD
		var k int

		// First, place the new value into a bin
		bins[i] = new
		i = (i + 1) % nSamples

		// Then iterate across the bins, getting a mean and a variance
		for k = 0; k < nSamples; k++ {
			x := bins[k]
			oldMean := Mean
			Mean = Mean + (x-Mean)/float64(k+1)
			S = S + (x-Mean)*(x-oldMean)
			//log.Printf("x = %g, mean = %g , k+1 = %d\n", x, Mean, k+1)
		}

		// return the mean and SD
		return Mean, math.Sqrt(S / float64(nSamples-1))
	}
}

// Mock generates a  mock "add" function, that computes a
// moving average and standard deviation of 0 and 1, respectively.
// Used to validate the math in tests
func Mock(nSamples int) func(s float64) (float64, float64) {
	return func(new float64) (float64, float64) {
		return 0.0, 1.0
	}
}
