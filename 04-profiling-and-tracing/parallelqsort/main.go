package main

import (
	"go-trace/qsort"
	"math/rand"
	"os"
	"runtime/trace"
)

func main() {
	// create random data to be sorted
	data := make([]int, 1e5)
	for i := range data {
		data[i] = rand.Int()
	}

	// start tracing program events to trace.out file
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := trace.Start(f); err != nil {
		panic(err)
	}
	defer trace.Stop()

	// sort data using parallel quicksort
	qsort.Sort(data)
}
