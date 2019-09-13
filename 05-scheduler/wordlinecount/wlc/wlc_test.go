package wlc

import (
	"io"
	"strings"
	"testing"
)

var testCases = []struct {
	descr    string
	haystack string
	needle   string
	want     uint64
}{
	{
		descr:    "empty haystack",
		haystack: "",
		needle:   "Go",
		want:     0,
	},
	{
		descr:    "empty needle",
		haystack: "Hello World",
		needle:   "",
		want:     1,
	},
	{
		descr:    "empty needle and haystack",
		haystack: "",
		needle:   "",
		want:     0,
	},
	{
		descr:    "single line, single occurrence",
		haystack: "Go",
		needle:   "Go",
		want:     1,
	},
	{
		descr: "no occurrences",
		haystack: `
Premature optimization is the root of all evil (or at least most of it) in programming.
Premature optimization is the root of all evil (or at least most of it) in programming.
Premature optimization is the root of all evil (or at least most of it) in programming.
Premature optimization is the root of all evil (or at least most of it) in programming.`,
		needle: "Knuth",
		want:   0,
	},
	{
		descr: "multiple occurrences",
		haystack: `
If police police police police, who polices the police police? Police police police police police police.
If police police police police, who polices the police police? Police police police police police police.
If police police police police, who polices the police police? Police police police police police police.
If police police police police, who polices the police police? Police police police police police police.`,
		needle: "police",
		want:   4,
	},
	{
		descr: "case insensitive",
		haystack: `
If police police police police, who polices the police police? Police police police police police police.
If police police police police, who polices the police police? Police police police police police police.
If police police police police, who polices the police police? Police police police police police police.
If police police police police, who polices the police police? Police police police police police police.`,
		needle: "PoLiCe",
		want:   4,
	},
}

func TestWordCountWordCountSeq(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.descr, func(t *testing.T) {
			got := WordLineCount([]io.Reader{strings.NewReader(tc.haystack)}, tc.needle)
			if got != tc.want {
				t.Fatalf("WordLineCountSeq(%q, %q)\ngot: %d\nwant: %d", tc.haystack, tc.needle, got, tc.want)
			}
		})
	}
}
