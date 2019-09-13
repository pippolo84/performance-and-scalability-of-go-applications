package wlc

import (
	"bufio"
	"io"
	"strings"
)

// WordLineCount counts the number of lines in all readers
// which contains an occurrence of the word w
func WordLineCount(readers []io.Reader, w string) uint64 {
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
