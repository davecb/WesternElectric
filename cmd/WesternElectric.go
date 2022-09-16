package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	movingAverage "github.com/davecb/WesternElectric/pkg/MovingAverage"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

/*
 * Western Electric Rules -- a classic set of decision rules in statistical process control
 * for detecting out-of-control or non-random conditions.
 */

func usage() {
	//nolint
	fmt.Fprint(os.Stderr, "Usage: westernelectric --samples N {file|-}\n") //nolint
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	var nSamples, mode int
	var report, table bool

	flag.IntVar(&nSamples, "nSamples", 5, "number of samples in the moving average")
	flag.BoolVar(&report, "report", false, "report anomalies only")
	flag.BoolVar(&table, "table", false, "report table of results & anomalies (default)")
	flag.Parse()

	switch {
	case report && table:
		log.Printf("Both table and report specified, choose only one. Halting\n")
		usage()
	case table:
		mode = 0
	case report:
		mode = 1
	}

	if flag.NArg() < 1 {
		fmt.Fprint(os.Stderr, "You must supply an input file, or '-' and a stream on stdin\n\n") //nolint
		usage()
	}
	if nSamples < 2 {
		fmt.Fprintf(os.Stderr, "You must specify a number of samples > 1 for the moving average, observed %d\n\n", nSamples) //nolint
		usage()
	}
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime) // show file:line in logs

	filename := flag.Arg(0)

	rc := WesternElectric(filename, nSamples, mode)
	os.Exit(rc)
}

// WesternElectric applies the W.E. rules to a stream of data, using a
// moving average of nSamples as the thing to compare against.
func WesternElectric(filename string, nSamples, mode int) int {
	var fp *os.File
	var err error

	if filename == "-" {
		// if the filename is "-", read stdin
		fp = os.Stdin
	} else {
		fp, err = os.Open(filename) //nolint
		if err != nil {
			log.Fatalf("error opening %s: %q, halting.", filename, err)
		}
		defer func() {
			err := fp.Close()
			if err != nil {
				log.Printf("Close of input file %q failed, ignored. %v\n",
					filename, err)
			}
		}()
	}
	rc := worker(fp, nSamples, mode)
	return rc
}

// worker reads the input and applies the rules, comparing the data
// to a moving average. For testing convenience, it returns the last error.
func worker(fp *os.File, nSamples, mode int) int { // FIXME return the string
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
			three := threeSigma(datum, average, sd)
			switch three {
			case " -3σ":
				lastErr = -3
			case " 3σ":
				lastErr = 3
			}
			two := twoSigma(datum, average, sd)
			switch two {
			case " -2σ":
				lastErr = -2
			case " 2σ":
				lastErr = 2
			}
			one := ""

			report(mode, date, datum, average, sd, three, two, one)
		}
		average, sd = add(datum)
	}
	return lastErr
}

func report(mode int, date string, datum float64, average float64, sd float64, three string, two string, one string) {
	// 	print stats and a visual indicator of broken rules

	switch mode {
	case 0: // print a table
		fmt.Printf("%s\t%f\t%0.4f\t%0.4f%s\t%s\t%s\n", date, datum, average, sd, three, two, one)

	case 1:
		// just a report
		fmt.Printf("%s\t%f\t%0.4f\t%0.4f%s\t%s\t%s\n", date, datum, average, sd, three, two, one)
	}
}

func header(mode int) {
	switch mode {
	case 0: // print a table
		fmt.Printf("%s\t                    %s\t%s\t    %s     %s\n", "#date", "datum", "average", "stddev", "flags")
	// 01-02T12:10:00-05:00
	case 1: // just a report
		fmt.Printf("%s\t                    %s\t%s\t    %s     %s\n", "#date", "datum", "average", "stddev", "flags")
	}
}

// threeSigma does the classic single-sample at 3 sigma test and returns a string
// to identify errors
func threeSigma(datum, average, sd float64) string {
	if math.Abs(datum) > average+(3*sd) {
		if datum > 0 {
			return " 3σ"
		} else {
			return " -3σ"
		}
	}
	return ""
}

// twoSigma detects 2 out of 3 points at +/- 2 sigma
func twoSigma(datum, average, sd float64) string {

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
			return " 2σ"
		} else {
			return " -2σ"
		}
	}
	threeSamples = shiftRight(threeSamples)
	return ""
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
