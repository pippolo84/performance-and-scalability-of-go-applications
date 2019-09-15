package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"performance-and-scalability-of-go-applications/05-scheduler/wordlinecount/solution/wlc"
	"runtime/trace"
)

// URL of the txt file to analyze
const url string = "https://www.gutenberg.org/files/16/16-0.txt"

// number of files to analyze
const nReaders int = 8

func main() {
	word := flag.String("word", "", "word to search")
	goroutines := flag.Int("goroutines", 0, "number of goroutines (0 for sequential)")
	traceFile := flag.String("trace", "", "trace file name")
	flag.Parse()

	// generate a trace file if requested
	if *traceFile != "" {
		f, err := os.Create(*traceFile)
		if err != nil {
			log.Fatalf("could not create trace output file %q: %v", *traceFile, err)
		}
		if err := trace.Start(f); err != nil {
			log.Fatalf("could not start trace: %v", err)
		}
		defer trace.Stop()
	}

	// in order to keep it simple, we read nReaders time from the same URL
	readers := make([]io.Reader, nReaders)
	for i := range readers {
		// Do the GET request
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("could not perform GET request to url %s", url)
		}
		defer resp.Body.Close()

		// save the http response body as reader interface
		readers[i] = resp.Body
	}

	// calculate total number of lines in which word appear,
	// using the requested number of goroutines
	var count uint64
	if *goroutines == 0 {
		count = wlc.WordLineCountSeq(readers, *word)
	} else {
		count = wlc.WordLineCountConc(readers, *word, *goroutines)
	}

	fmt.Println(count)
}
