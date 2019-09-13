package wlc

import (
	"bufio"
	"io"
	"strings"
	"sync"
	"sync/atomic"
)

// WordLineCountSeq counts (sequentially) the number of lines in all readers
// which contains an occurrence of the word w
func WordLineCountSeq(readers []io.Reader, w string) uint64 {
	// total counter
	var count uint64

	// store w as lowercase (case insensitive match)
	w = strings.ToLower(w)

	// for each reader
	for _, r := range readers {
		// set the scanner
		scanner := bufio.NewScanner(r)

		// Set the split function for the scanning operation.
		scanner.Split(bufio.ScanLines)

		// search the word in every line and count the lines where we found it
		for scanner.Scan() {
			if strings.Contains(strings.ToLower(scanner.Text()), w) {
				count++
			}
		}
	}

	// return results
	return count
}

// WordLineCountConc counts (concurrently) the number of lines in all readers
// which contains an occurrence of the word wm using nGoroutines
func WordLineCountConc(readers []io.Reader, w string, nGoroutines int) uint64 {
	// total counter
	var count uint64

	// store w as lowercase
	w = strings.ToLower(w)

	// create a buffered channel to send readers to goroutines
	ch := make(chan io.Reader, len(readers))
	for _, r := range readers {
		ch <- r
	}
	close(ch)

	var wg sync.WaitGroup
	wg.Add(nGoroutines)

	// create goroutines to crunch data
	for i := 0; i < nGoroutines; i++ {
		go func() {
			defer wg.Done()

			// loop over the channel to get the data source
			for r := range ch {
				// set the scanner
				scanner := bufio.NewScanner(r)

				// Set the split function for the scanning operation.
				scanner.Split(bufio.ScanLines)

				// local goroutine counter
				var gCount uint64

				// search for word w in line l
				for scanner.Scan() {
					if strings.Contains(strings.ToLower(scanner.Text()), w) {
						gCount++
					}
				}

				// add local count to the global one
				atomic.AddUint64(&count, gCount)
			}
		}()
	}

	// wait for the goroutines to complete work
	wg.Wait()

	// return results
	return count
}
