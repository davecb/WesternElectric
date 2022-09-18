package main

import (
	"flag"
	"fmt"
	we "github.com/davecb/WesternElectric/cmd/WesternElectric"
	"log"
	"os"
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
	var nSamples, reportingMode int
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
		reportingMode = 0
	case report:
		reportingMode = 1
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

	rc := we.ApplyRules(filename, nSamples, reportingMode)
	os.Exit(rc)
}
