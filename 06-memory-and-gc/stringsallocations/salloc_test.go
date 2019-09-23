package salloc

import (
	"testing"
	"time"
)

type strcatInput struct {
	id   int
	name string
	date time.Time
}

var testCases = []struct {
	descr string
	input strcatInput
	want  string
}{
	{
		descr: "simple test",
		input: strcatInput{
			id:   9,
			name: "Go workshop",
			date: time.Date(2019, time.September, 24, 16, 30, 0, 0, time.UTC),
		},
		want: "9 Go workshop 2019-09-24 16:30:00 +0000 UTC",
	},
}

func TestStrcat(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.descr, func(t *testing.T) {
			got := Strcat(tc.input.id, tc.input.name, tc.input.date)
			if got != tc.want {
				t.Fatalf("Strcat(...) = %q, want %q", got, tc.want)
			}
		})
	}
}
