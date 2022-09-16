package WesternElectric

import (
	"log"
	"os"
)

// ApplyRules applies the Western Electric rules to a stream of data, using a
// moving average of nSamples as the thing to compare against.
func ApplyRules(filename string, nSamples, mode int) int {
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
	rc := Worker(fp, nSamples, mode)
	return rc
}
