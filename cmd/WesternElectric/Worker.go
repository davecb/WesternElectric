package WesternElectric

import (
	"encoding/csv"
	"fmt"
	movingAverage "github.com/davecb/WesternElectric/pkg/MovingAverage"

	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

// worker reads the input and applies the rules, comparing the data
// to a moving average. For testing convenience, it returns the last anomaly.
func Worker(fp *os.File, nSamples, mode int) int {
	var nr, lastErr int
	var average float64
	var sd float64

	// set up csv reader to read fields out of a file
	r := csv.NewReader(fp)
	r.Comma = ' '
	r.Comment = '#'
	r.FieldsPerRecord = -1 // ignore differences
	r.LazyQuotes = true    // allow bad quoting

	// set up moving average
	add := movingAverage.New(nSamples)

	// read lines containing a datestamp or other initial field, and a value
	header(mode)
	for nr = 0; ; nr++ {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// we had a csv-reading error, die.
			log.Fatalf("error %q, in %q, line %d\n", err, record, nr)
		}
		if len(record) < 2 {
			// skip it, but complain
			log.Printf("Too few fields in line %d, %q. Ignored.\n", nr, record)
			continue
		}
		//log.Printf("read %q\n", record)

		date := record[0]
		// parse the value field
		datum, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			// we had a float-parsing error
			log.Printf("Invalid float64 in line %d, %q. Ignored.\n", nr, strings.Join(record, "\t"))
			continue
		}

		//log.Printf("at time %q, got %g, average = %g, sd = %g\n", record[0], datum, average, sd)
		if nr > nSamples {
			// see if we break any of the rules, but only once we have an average to use
			rcThree := ThreeSigma(datum, average, sd)
			if rcThree != 0 {
				lastErr = rcThree
			}

			rcTwo := TwoSigma(datum, average, sd)
			if rcTwo != 0 {
				lastErr = rcTwo
			}
			rcOne := 0

			report(mode, date, datum, average, sd, rcThree, rcTwo, rcOne)
		}
		average, sd = add(datum)
	}
	return lastErr
}

// report tells us what happened, in short or long form.
func report(mode int, date string, datum float64, average float64, sd float64, rcThree, rcTwo, rcOne int) {
	// 	print stats and a visual indicator of broken rules
	var three, two, one string

	switch rcThree {
	case -3:
		three = " -3σ"
	case 3:
		three = " 3σ"
	default:
		three = ""
	}
	switch rcTwo {
	case -3:
		two = " -2σ"
	case 3:
		two = " 2σ"
	default:
		two = ""
	}

	switch mode {
	case 0: // print a table of date, datum and the +/- sigma lines, then the indicators
		fmt.Printf("%s %0.4f %0.4f %0.4f %0.4f %0.4f %0.4f %0.4f %0.4f %s %s %s\n",
			date, datum, average,
			average+sd, average-sd,
			average+2*sd, average-2*sd,
			average+3*sd, average-3*sd,
			three, two, one)

	case 1:
		// just a report, for people to read
		fmt.Printf("%s %f %0.4f %0.4f %s %s %s\n", date, datum, average, sd, three, two, one)
	}
}

func header(mode int) {
	switch mode {
	case 0: // print a table for plotting and/or spreadsheets
		fmt.Printf("#date datum average average+sd average-sd average+2*sd average-2*sd average+3*sd average-3*sd flags\n")
	case 1: // just a report, aligned for people to scan
		fmt.Printf("%s %s         %s     %s      %s\n", "#date", "datum", "average", "stddev", "flags")
	}
}

// ThreeSigma does the classic single-sample at 3 sigma test and returns a string
// to identify errors
func ThreeSigma(datum, average, sd float64) int {
	if math.Abs(datum) > average+(3*sd) {
		if datum > 0 {
			return 3
		} else {
			return -3
		}
	}
	return 0
}

// TwoSigma detects 2 out of 3 points at +/- 2 sigma
func TwoSigma(datum, average, sd float64) int {

	// record its state
	switch {
	case datum > average+(2*sd):
		threeSamples[0] = StateAbove
	case datum < average-(2*sd):
		threeSamples[0] = StateBelow
	default:
		threeSamples[0] = StateNA
	}
	if twoOf(threeSamples) {
		threeSamples = shiftRight(threeSamples)
		if datum > 0 {
			return 2
		} else {
			return -2
		}
	}
	threeSamples = shiftRight(threeSamples)
	return 0
}

// oneSigma detects  4/5 at 1 +/- sigma

// noSigma detects  9/9 on the same side of 0

/*
 * infrastructure for the tests
 */
var threeSamples, fiveSamples []State

func init() {
	// state vectors for twp of three, four of five
	threeSamples = make([]State, 3)
	fiveSamples = make([]State, 5)
}

// twoOf reports true if two states match
func twoOf(twosies []State) bool {
	return nOf(twosies, 2)
}

// nOf reports true if N states match, including the first
func nOf(window []State, matches int) bool {
	var found int
	var target = window[0]

	if target == StateNA {
		// we only care about above or belows matching
		return false
	}

	for i := 0; i < len(window); i++ {
		if window[i] == target {
			found++
		}
	}
	if found >= matches {
		return true
	}
	return false
}

// shiftRight moves everything to the right, zero-filling
func shiftRight(window []State) []State {
	for i := len(window) - 1; i > 0; i-- {
		window[i] = window[i-1]
	}
	window[0] = StateNA
	return window
}

// State is the state of a previous sample. NA means it was neither
// above nor below the +/- cutoff, //Above means it was above the cutoff, and so on
type State int32

const (
	StateNA    State = 0
	StateAbove State = 1
	StateBelow State = 2
)

var StateName = map[int32]string{
	0: "NA",
	1: "Above",
	2: "Below",
}

func (x State) String() string {
	return StateName[int32(x)]
}
