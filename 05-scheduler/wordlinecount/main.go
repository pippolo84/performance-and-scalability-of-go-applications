package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"performance-and-scalability-of-go-applications/05-scheduler/wordlinecount/wlc"
)

// URL of the txt file to analyze
const url string = "https://www.gutenberg.org/files/16/16-0.txt"

// number of files to analyze
const nReaders int = 8

func main() {
	word := flag.String("word", "", "word to search")
	flag.Parse()

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

	// calculate total number of lines in which word appear
	fmt.Println(wlc.WordLineCount(readers, *word))
}
