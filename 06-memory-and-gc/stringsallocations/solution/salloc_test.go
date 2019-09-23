package salloc

import (
	"testing"
	"time"
)

type impl struct {
	descr string
	f     func(id int, name string, t time.Time) string
}

var impls = []impl{
	{"StrcatStrAppend", StrcatStrAppend},
	{"StrcatStrBuilder", StrcatStrBuilder},
	{"StrcatSprintf", StrcatSprintf},
	{"StrcatBytesBuffer", StrcatBytesBuffer},
	{"StrcatSlicePreallocated", StrcatSlicePreallocated},
	{"StrcatByteSlice", StrcatByteSlice},
}

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
	for _, imp := range impls {
		t.Run(imp.descr, func(t *testing.T) {
			for _, tc := range testCases {
				t.Run(tc.descr, func(t *testing.T) {
					got := imp.f(tc.input.id, tc.input.name, tc.input.date)
					if got != tc.want {
						t.Fatalf("Strcat(...) = %q, want %q", got, tc.want)
					}
				})
			}
		})
	}
}

func BenchmarkStrcat(b *testing.B) {
	for _, imp := range impls {
		b.Run(imp.descr, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = imp.f(testCases[0].input.id, testCases[0].input.name, testCases[0].input.date)
			}
		})
	}
}
