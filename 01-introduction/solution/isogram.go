package isogram

import (
	"sort"
	"strings"
	"unicode"
)

// IsIsogram implemented using a map
func IsIsogramMap(s string) bool {
	found := make(map[rune]bool)

	for _, r := range s {
		if !unicode.IsLetter(r) {
			continue
		}

		r = unicode.ToLower(r)
		if found[r] {
			return false
		}

		found[r] = true
	}

	return true
}

// IsIsogram implemented using sorting
func IsIsogramSorting(s string) bool {
	if len(s) < 2 {
		return true
	}

	runes := []rune(strings.ToLower(s))
	sort.Sort(runeSlice(runes))

	for i := 1; i < len(runes); i++ {
		if runes[i-1] == runes[i] && unicode.IsLetter(runes[i-1]) {
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

// IsIsogram implemented using linear search
func IsIsogramLinearSearch(s string) bool {
	s = strings.ToLower(s)

	for i, c := range s {
		if unicode.IsLetter(c) && strings.ContainsRune(s[i+1:], c) {
			return false
		}
	}

	return true
}

// IsIsogram implemented using a bool array
func IsIsogramBoolArray(s string) bool {
	foundRune := [26]bool{} //'a' to 'z'

	for _, r := range s {
		if !unicode.IsLetter(r) {
			continue
		}

		// convert the rune to lowercase to index foundRune
		r = unicode.ToLower(r)
		i := r - 'a'

		if foundRune[i] == true {
			return false
		}
		foundRune[i] = true
	}

	return true
}
