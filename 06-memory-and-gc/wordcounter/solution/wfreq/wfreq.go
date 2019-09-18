package wfreq

import (
	"sort"
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
		start := 0
		for start < len(text) && (text[start] == ' ' || text[start] == '\n') {
			start++
		}
		i := start + 1
		for i < len(text) {
			if text[start] == ' ' || text[start] == '\n' {
				start++
			}
			if (text[i] == ' ' || text[i] == '\n') && start < i {
				wfreq[text[start:i]]++
				start = i
			}
			i++
		}
		if start < len(text) && (text[start] != ' ' && text[start] != '\n') {
			wfreq[text[start:len(text)]]++
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
