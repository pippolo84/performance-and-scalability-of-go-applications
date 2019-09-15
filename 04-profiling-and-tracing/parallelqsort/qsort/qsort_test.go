package qsort

import (
	"reflect"
	"testing"
)

var testCases = []struct {
	descr string
	input []int
	want  []int
}{
	{
		descr: "empty slice",
		input: []int{},
		want:  []int{},
	},
	{
		descr: "single element slice",
		input: []int{10},
		want:  []int{10},
	},
	{
		descr: "ordered two elements slice",
		input: []int{1, 10},
		want:  []int{1, 10},
	},
	{
		descr: "unordered two elements slice",
		input: []int{10, -15},
		want:  []int{-15, 10},
	},
	{
		descr: "unordered slice",
		input: []int{10, -15, -5, 17, 24, 4, 2, -36, 103},
		want:  []int{-36, -15, -5, 2, 4, 10, 17, 24, 103},
	},
}

func TestInsertRemove(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.descr, func(t *testing.T) {
			Sort(tc.input)
			if !reflect.DeepEqual(tc.input, tc.want) {
				t.Fatalf("got %v, want %v", tc.input, tc.want)
			}
		})
	}
}
