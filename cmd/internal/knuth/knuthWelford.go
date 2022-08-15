package main

import (
	"fmt"
	"math"
)

// see https://pkg.go.dev/github.com/eclesh/welford#section-readme -- does infinite sequence in go

//https://www.johndcook.com/blog/2008/09/26/comparing-three-methods-of-computing-standard-deviation/
//https://www.johndcook.com/blog/standard_deviation/

//https://stackoverflow.com/questions/1174984/how-to-efficiently-calculate-a-running-standard-deviation

func create(nSamples int) func(s float64) (float64, float64) {
	var Mean, S float64 // 0.0
	var i, k int

	bins := make([]float64, nSamples)
	return func(new float64) (float64, float64) {
		// First place the new value into the bin
		bins[i] = new
		i = (i + 1) % nSamples

		// Then iterate across the bins, getting a mean and a variance
		for k = 0; k < nSamples; k++ {
			x := bins[k]
			oldMean := Mean
			Mean = Mean + (x-Mean)/float64(k+1)
			S = S + (x-Mean)*(x-oldMean)
		}

		return Mean, math.Sqrt(S / float64(nSamples-1))
	}
}

func main() {
	add := create(5)
	add(1)
	add(2)
	add(3)
	add(4)
	a, _ := add(5)
	fmt.Println("(1+2+3+4+5          ) / 5 =", a)
	a, _ = add(9)
	fmt.Println("(  2+3+4+5+9        ) / 5 =", a)
	a, _ = add(3)
	fmt.Println("(    3+4+5+9+3      ) / 5 =", a)
	a, _ = add(0)
	fmt.Println("(      4+5+9+3+0    ) / 5 =", a)
	a, _ = add(-9)
	fmt.Println("(        5+9+3+0-9  ) / 5 =", a)
	a, _ = add(-8)
	fmt.Println("(          9+3+0-9-8) / 5 =", a)
}

//https://medium.com/@tobMTV/19-streaming-algorithm-for-the-mean-and-the-variance-b8a9d8faf8a6
// Tobia Vendrame
//variance(samples):
//M := 0
//S := 0
//for k from 1 to N:
//	x := samples[k]
//	oldM := M
//	M := M + (x-M)/k
//	S := S + (x-M)*(x-oldM)
//return S/(N-1)
