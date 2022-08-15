package movingAverage

import "fmt"

/*
 * start with 5 slots
 * add 1 through 5
 * get avg, sd
 * add some known values and compare
 * do the diagrammatic style, and cite the author

 * continue with a new 5 slots
 * add nothing
 * get nothing
 * add four
 * get the average and sd, which should still work
 */

func main() {
	add := NewFMovingAverage(5)
	add(1)
	add(2)
	add(3)
	add(4)
	fmt.Println("(1+2+3+4+5          ) / 5 =", add(5))
	fmt.Println("(  2+3+4+5+9        ) / 5 =", add(9))
	fmt.Println("(    3+4+5+9+3      ) / 5 =", add(3))
	fmt.Println("(      4+5+9+3+0    ) / 5 =", add(0))
	fmt.Println("(        5+9+3+0-9  ) / 5 =", add(-9))
	fmt.Println("(          9+3+0-9-8) / 5 =", add(-8))
}

//$ go run roll.go
//(1+2+3+4+5          ) / 5 = 3
//(  2+3+4+5+9        ) / 5 = 4.6
//(    3+4+5+9+3      ) / 5 = 4.8
//(      4+5+9+3+0    ) / 5 = 4.2
//(        5+9+3+0-9  ) / 5 = 1.6
//(          9+3+0-9-8) / 5 = -1
