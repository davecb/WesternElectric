package WesternElectric

import (
	WesternElectric2 "github.com/davecb/WesternElectric/cmd/WesternElectric"
	movingAverage "github.com/davecb/WesternElectric/pkg/MovingAverage"
	"log"
	"testing"
)

func Test_threeSigma(t *testing.T) {
	// threeSigma does the classic single-sample test
	tests := []struct {
		name     string
		data     []float64 // last value must be the outlier
		nSamples int       // number of samples in moving average, must be > 1
		expect   string
	}{
		{
			name:     "Not 3 sigma",
			data:     []float64{1, 2, 3, 4, 5, 9, 3, 0, 1},
			nSamples: 5,
			expect:   "",
		},
		{
			name:     "3 sigma",
			data:     []float64{1, 2, 3, 4, 5, 9, 3, 0, 99},
			nSamples: 5,
			expect:   " 3σ",
		},
		{
			name:     "3 sigma, smaller",
			data:     []float64{1, 2, 3, 4, 5, 9, 3, 0, 22},
			nSamples: 5,
			expect:   " 3σ",
		},
		{
			name:     "-3 sigma",
			data:     []float64{1, 2, 3, 4, 5, 9, 3, 0, -99},
			nSamples: 5,
			expect:   " -3σ",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			var datum, average, sd float64
			var oldAv, oldSD float64
			var got string
			var i int

			t.Logf("n:  flag  data oldAv  oldSD  a+3sd\n")
			add := movingAverage.New(tt.nSamples)
			for i, datum = range tt.data {
				got = WesternElectric2.threeSigma(datum, average, sd)
				average, sd = add(datum)
				oldAv, oldSD = average, sd
			}
			t.Logf("%-0.2d: %-5.5q %-4.2g %-4.2g %-4.2g %-4.2g\n", i, got, datum, oldAv, oldSD, oldAv+3*oldSD)
			if got != tt.expect {
				t.Errorf("threeSigma() = %q, expected %q", got, tt.expect)
			}

		})
	}
}

func Test_twoSigma(t *testing.T) {
	// twoSigma detects 2 out of 3 points at +/- 2 sigma
	// for the tests, the last entry MUST be the one that reports, that's
	// the only one we check.
	tests := []struct {
		name     string
		data     []float64 // last value must be the outlier
		nSamples int       // number of samples in moving average, must be >= 1 unless mocked
		mock     bool
		expect   string
	}{
		{
			name:     "2-sigma",
			data:     []float64{1, 1, 1, 1, 1, 1, 5, 1, 11},
			nSamples: 5,
			expect:   " 2σ",
		},
		{
			name:     "2-sigma mock 1", // , filler, then last two out of three
			data:     []float64{1, 1, 4, 4},
			nSamples: 1, // OK iff mocked
			mock:     true,
			expect:   " 2σ",
		},
		{
			name:     "2-sigma mock 2",
			data:     []float64{1, 1, 4, 1, 4}, // filler, then first and last out of three
			nSamples: 1,
			mock:     true,
			expect:   " 2σ",
		},
		{
			name:     "2-sigma mock 3",
			data:     []float64{1, 1, 4, 4}, // filler, then first and second out of, well, two
			nSamples: 1,
			mock:     true,
			expect:   " 2σ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var datum, average, sd float64
			var add func(s float64) (float64, float64)
			var oldAv, oldSD float64
			var got string
			var i int

			if tt.mock {
				// use average == 0, sd == 1
				add = movingAverage.Mock(tt.nSamples)
			} else {
				add = movingAverage.New(tt.nSamples)
			}
			t.Logf("n:  flag  data av   sd   a+3sd\n")
			for i, datum = range tt.data {
				if i > tt.nSamples {
					// once we have an average, start looking
					got = WesternElectric2.twoSigma(datum, average, sd)
					t.Logf("%-0.2d: %-5.5q %-4.2g %-4.2g %-4.2g %-4.2g\n", i, got, datum, oldAv, oldSD, oldAv+3*oldSD)
					oldAv, oldSD = average, sd
				}
				average, sd = add(datum)
			}
			//t.Logf("%-0.2d: %-5.5q %-4.2g %-4.2g %-4.2g %-4.2g\n", i, got, datum, oldAv, oldSD, oldAv+3*oldSD)
			if got != tt.expect {
				t.Errorf("twoSigma() = %q, expected %q", got, tt.expect)
			}

		})
	}
}

// TestShiftRight looks to see if shiftRight(1,2,3) = 0,1,2
// where the leftmost value is zero-filled.
func Test_shiftRight(t *testing.T) {
	var vector [3]WesternElectric2.State

	vector[0] = 0
	vector[1] = 1
	vector[2] = 2
	x := WesternElectric2.shiftRight(vector[:])
	if x[0] != 0 || x[1] != 0 || x[2] != 1 {
		t.Errorf("shiftRight failed\n")
	}
}

// Test WE with files of data
func Test_westernElectric(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		nSamples int // number to average, not number samples in data
		expect   int // a sigma indication
	}{
		{
			name:     "example.csv", // FIXME, breake up into sets
			file:     "./testdata/example.csv",
			nSamples: 5,
			expect:   2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := WesternElectric(tt.file, tt.nSamples)
			if rc != tt.expect {
				t.Errorf("we found a failure\n")
			}

		})
	}
}

func ExampleWesternElectric() {
	rc := WesternElectric("./testdata/example.csv", 5)
	if rc > 0 {
		log.Printf("we found at least one failure\n")
	}

}
