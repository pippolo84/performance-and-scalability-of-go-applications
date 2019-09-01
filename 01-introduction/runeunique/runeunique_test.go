package runeunique

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
	"unicode/utf8"
)

type testCase struct {
	descr string
	input string
}

var testCases []testCase

var implementations = []struct {
	descr string
	f     func(s string) bool
}{
	{
		descr: "01-map",
		f:     IsRuneUniqueMap,
	},
	{
		descr: "02-sorting",
		f:     IsRuneUniqueSorting,
	},
	{
		descr: "03-linear-search",
		f:     IsRuneUniqueLinearSearch,
	},
}

func generateRunes(n int) []rune {
	runes := make([]rune, n)

	i, j := 0, 0
	for j < n {
		r := rune(i)
		if utf8.ValidRune(r) {
			runes[j] = r
			j++
		}
		i++
	}

	return runes
}

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	for _, n := range []int{8, 16, 32, 64, 128, 256, 512, 1024, 4096, 8192, 16384} {
		runes := generateRunes(n)
		// shuffle runes to make the sorting implementation acting fair
		rand.Shuffle(len(runes), func(i, j int) { runes[i], runes[j] = runes[j], runes[i] })
		testCases = append(testCases, testCase{fmt.Sprintf("%d", n), string(runes)})
	}

	os.Exit(m.Run())
}

func TestIsRuneUnique(t *testing.T) {
	for _, impl := range implementations {
		t.Run(impl.descr, func(t *testing.T) {
			for _, tc := range testCases {
				t.Run(tc.descr, func(t *testing.T) {
					got := impl.f(tc.input)
					if got != true {
						t.Fatalf("IsRuneUnique(%q)\ngot: %t\nwant: %t", tc.input, got, true)
					}
				})
			}
		})
	}
}

func BenchmarkIsIsogram(b *testing.B) {
	for _, impl := range implementations {
		b.Run(impl.descr, func(b *testing.B) {
			for _, tc := range testCases {
				b.Run(tc.descr, func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						impl.f(tc.input)
					}
				})
			}
		})
	}
}
