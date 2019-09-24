package main

import (
	"os"
	"runtime"
	"runtime/pprof"
)

type s struct {
	v int
}

func newStruct() *s {
	return &s{10}
}

var globX []*s

func f(l int) []*s {
	var x []*s

	for i := 0; i < l; i++ {
		x = append(x, newStruct())
	}

	return x
}

func main() {
	runtime.MemProfileRate = 1

	// does not allocate
	x := newStruct()

	println(x.v)

	// allocates on the heap
	globX = f(5)

	f, err := os.Create("heap.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		panic(err)
	}
}
