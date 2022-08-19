package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	movingAverage "github.com/davecb/WesternElectric/pkg/MovingAverage"
	"io"
	"log"
	"os"
	"strconv"
)

/*
 * Western Electric Rules -- a classic set of decision rules in statistical process control
 * for detecting out-of-control or non-random conditions.
 */

func usage() {
	//nolint
	fmt.Fprint(os.Stderr, "Usage: westernelectrictc --samples N {file|-}\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	var nSamples int
	var fp *os.File
	var err error

	flag.IntVar(&nSamples, "nSamples", 5, "number of nSamples in the moving average")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprint(os.Stderr, "You must supply an input file or '-' and a stream on stdin\n\n") //nolint
		usage()
	}
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime) // show file:line in logs

	filename := flag.Arg(0)
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
	worker(fp, nSamples)
}

func worker(fp *os.File, nSamples int) {
	var nr int
	var average float64
	var sd float64

	// set up csv reader to read fields out of a file
	r := csv.NewReader(fp)
	r.Comma = '\t'
	r.Comment = '#'
	r.FieldsPerRecord = -1 // ignore differences
	r.LazyQuotes = true    // allow bad quoting

	// set up moving average
	add := movingAverage.CreateMovingAverage(nSamples)

	// read lines containing a data and a (usually floating-point) value
	for nr = 0; ; nr++ {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("error %q, in %q, line %d\n", err, record, nr)
		}
		datum, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatalf("Encountered invalid float64 in %q, one line %d, halting.\n", record, nr)
		}
		average, sd = add(datum)
		log.Printf("got %f, average = %f, sd = %f\n", datum, average, sd)
		// 	add
		//  test
		// 	print stats
	}
	log.Printf("out of loop at line %d\n", nr)
}
