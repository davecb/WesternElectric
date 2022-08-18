package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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
	var samples int
	var fp *os.File

	flag.IntVar(&samples, "samples", 5, "number of samples in the moving average")
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
		f, err := os.Open(filename) //nolint
		if err != nil {
			log.Fatalf("error opening %s: %q, halting.", filename, err)
		}
		defer func() {
			err := f.Close()
			if err != nil {
				log.Printf("Close of r/o input file %q failed, ignored. %v",
					filename, err)
			}
		}()
	}
	worker(fp, samples)
	// create a worker
	// run it
	return
}

func worker(f *os.File, samples int) {
	// while read stdin {
	//		add
	//		test
	// 		print stats
	// }
	return

}
