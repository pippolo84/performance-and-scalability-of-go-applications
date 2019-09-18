package wfreq

import (
	"sort"
	"strings"
)

type (
	Counter struct {
		Word      string `json:"word"`
		Frequency int    `json:"frequency"`
	}

	UserWords struct {
		ID    int       `json:"id"`
		Name  string    `json:"name"`
		Words []Counter `json:"words"`
	}
)

// WordsFreq returns a slice containing the frequencies of each word in texts
// the slice is sorted in descending order of word frequencies
func WordsFreq(texts []string) []Counter {
	// split all the text lines in words, counting frequencies
	wfreq := make(map[string]int)
	for _, text := range texts {
		for _, word := range strings.Fields(text) {
			wfreq[word]++
		}
	}

	// sort words by frequencies
	var wcounters []Counter
	for word, freq := range wfreq {
		wcounters = append(wcounters, Counter{word, freq})
	}
	sort.Slice(wcounters, func(i, j int) bool {
		if wcounters[i].Frequency == wcounters[j].Frequency {
			return wcounters[i].Word < wcounters[j].Word
		}
		return wcounters[i].Frequency > wcounters[j].Frequency
	})

	return wcounters
}
