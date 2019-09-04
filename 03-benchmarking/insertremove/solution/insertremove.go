// Package insertremove contains the Insert and Remove exercise implementations
// "Insert a sequence of random integers into a sorted sequence, then remove those
// elements one by one as determined by a random sequence of positions"
package insertremove

import (
	"errors"
	"math/rand"
	"sort"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// InsertRemoveSlice takes an integer and solves the problem for n elements using a slice
func InsertRemoveSlice(n int) error {
	const maxInt = int(^uint(0) >> 1)

	if n < 0 {
		return errors.New("n must be a non negative value")
	}

	values := make([]int, n)
	for i := range values {
		values[i] = maxInt
	}

	// insert
	for i := 0; i < n; i++ {
		cur := rand.Int()
		pos := sort.SearchInts(values, cur)
		copy(values[pos+1:], values[pos:])
		values[pos] = cur
	}

	// remove
	for n > 0 {
		pos := rand.Intn(n)
		values = append(values[:pos], values[pos+1:]...)
		n--
	}

	return nil
}

// InsertRemoveList takes an integer and solves the problem for n elements using a slice
func InsertRemoveList(n int) error {
	return nil
}
