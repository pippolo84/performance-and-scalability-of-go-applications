package runeunique

import (
	"sort"
	"strings"
)

// IsRuneUniqueMap implemented using a map
func IsRuneUniqueMap(s string) bool {
	found := make(map[rune]bool)

	for _, r := range s {
		if found[r] {
			return false
		}

		found[r] = true
	}

	return true
}

// IsRuneUniqueSorting implemented using sorting
func IsRuneUniqueSorting(s string) bool {
	if len(s) < 2 {
		return true
	}

	runes := []rune(s)
	sort.Sort(runeSlice(runes))

	for i := 1; i < len(runes); i++ {
		if runes[i-1] == runes[i] {
			return false
		}
	}

	return true
}

type runeSlice []rune

func (p runeSlice) Len() int {
	return len(p)
}

func (p runeSlice) Less(i, j int) bool {
	return p[i] < p[j]
}

func (p runeSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// IsRuneUniqueLinearSearch implemented using linear search
func IsRuneUniqueLinearSearch(s string) bool {
	for i, c := range s {
		if strings.ContainsRune(s[i+1:], c) {
			return false
		}
	}

	return true
}
