package qsort

import "sync"

// Sort sorts data slice using a parallel implementation of the quicksort algorithm
func Sort(data []int) {
	var wg sync.WaitGroup

	wg.Add(1)
	qSortPar(data, &wg)

	wg.Wait()
}

// data length threshold for using parallel sorting
const threshold = 1e3

func partition(data []int) int {
	// get the pivot
	pivot := data[0]

	// all values <= than pivot should pile up to the left,
	// while the others should pile up to the right
	left, right := 1, len(data)-1
	for right >= left {
		if data[left] <= pivot {
			left++
		} else {
			data[right], data[left] = data[left], data[right]
			right--
		}
	}

	// swap pivot into middle
	data[left-1], data[0] = data[0], data[left-1]

	return left
}

func qSortSeq(data []int) {
	// recursion base case
	if len(data) < 2 {
		return
	}

	// call a round of partition
	left := partition(data)

	// recursively sort subsets
	qSortSeq(data[:left-1])
	qSortSeq(data[left:])
}

func qSortPar(data []int, wg *sync.WaitGroup) {
	defer wg.Done()

	if len(data) < 2 {
		// we should end up here only if we make an initial call
		// to qSortPar with len(data) < 2
		return
	}

	// call a round of partition
	left := partition(data)

	// launch separate goroutines for partitions longer than threshold,
	// use the same goroutine to sort shorter partitions
	if left-1 > threshold {
		wg.Add(1)
		go qSortPar(data[:left-1], wg)
	} else {
		qSortSeq(data[:left-1])
	}

	if len(data)-left > threshold {
		wg.Add(1)
		go qSortPar(data[left:], wg)
	} else {
		qSortSeq(data[left:])
	}

}
