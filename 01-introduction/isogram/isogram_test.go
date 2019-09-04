package isogram

import "testing"

var testCases = []struct {
	description string
	input       string
	want        bool
}{
	{
		description: "empty string",
		input:       "",
		want:        true,
	},
	{
		description: "isogram with only lower case characters",
		input:       "isogram",
		want:        true,
	},
	{
		description: "word with one duplicated character",
		input:       "eleven",
		want:        false,
	},
	{
		description: "word with one duplicated character from the end of the alphabet",
		input:       "zzyzx",
		want:        false,
	},
	{
		description: "longest reported english isogram",
		input:       "subdermatoglyphic",
		want:        true,
	},
	{
		description: "word with duplicated character in mixed case",
		input:       "Alphabet",
		want:        false,
	},
	{
		description: "word with duplicated character in mixed case, lowercase first",
		input:       "alphAbet",
		want:        false,
	},
	{
		description: "hypothetical isogrammic word with hyphen",
		input:       "thumbscrew-japingly",
		want:        true,
	},
	{
		description: "isogram with duplicated hyphen",
		input:       "six-year-old",
		want:        true,
	},
	{
		description: "made-up name that is an isogram",
		input:       "Emily Jung Schwartzkopf",
		want:        true,
	},
	{
		description: "duplicated character in the middle",
		input:       "accentor",
		want:        false,
	},
	{
		description: "same first and last characters",
		input:       "angola",
		want:        false,
	},
}

func TestIsIsogram(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			got := IsIsogram(tc.input)
			if got != tc.want {
				t.Fatalf("IsIsogram(%q)\ngot: %t\nwant: %t", tc.input, got, tc.want)
			}
		})
	}
}

func BenchmarkIsIsogram(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, c := range testCases {
			IsIsogram(c.input)
		}
	}
}
