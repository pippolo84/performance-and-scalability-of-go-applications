package main

import (
	"os"
	"runtime"
	"runtime/trace"
	"strconv"
	"sync"
)

// number crunching in a tight loop
func tightLoop() {
	for j := 0; j < 100; j++ {
		a := 100
		for i := 1; i < 100000; i++ {
			a += i * 100 / i
		}
		for i := 100000/a + 1; i < 100000; i++ {
			a += i * 100 / i
		}
	}
}

// a lot of allocations to trigger the GC
func heavyAlloc() {
	var x []string

	for i := 0; i < 10000; i++ {
		x = append(x, "Hello "+strconv.Itoa(i)+" World")
	}
}

func main() {
	//start tracing program events to trace.out file
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := trace.Start(f); err != nil {
		panic(err)
	}
	defer trace.Stop()

	// start some goroutines to execute heavyAlloc
	// and some goroutines to execute tightLoop
	nGoroutines := runtime.GOMAXPROCS(0) * 2

	var wg sync.WaitGroup
	wg.Add(nGoroutines)
	for i := 0; i < nGoroutines/2; i++ {
		go func() {
			defer wg.Done()
			heavyAlloc()
		}()
		go func() {
			defer wg.Done()
			tightLoop()
		}()
	}

	wg.Wait()
}
