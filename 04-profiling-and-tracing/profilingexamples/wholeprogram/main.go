// Starts this program with
//
// go run main.go
//
// In other tab, collect the profile with
//
// go tool pprof http://localhost:9090/debug/pprof/profile

package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
)

// Fibonacci returns the n-th element of the Fibonacci sequence
func Fibonacci(n int) int {
	if n == 0 || n == 1 {
		return 1
	}

	return Fibonacci(n-2) + Fibonacci(n-1)
}

func main() {
	// add profiling and tracing flags
	traceFile := flag.String("trace", "", "write program trace to file")
	cpuProfile := flag.String("cpuprofile", "", "write cpu profile to file")
	memProfile := flag.String("memprofile", "", "write heap profile to file")
	flag.Parse()

	// start tracing program events to file
	if *traceFile != "" {
		f, err := os.Create(*traceFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if err := trace.Start(f); err != nil {
			panic(err)
		}
		defer trace.Stop()
	}

	// start CPU profiling and save it to file
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			panic(err)
		}

		if err := pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

	// set memory allocations sampling rate to the maximum
	// as soon as we are interested to profile allocations
	runtime.MemProfileRate = 1

	// code to profile
	fmt.Println(Fibonacci(42))

	// save memory allocations profile to file
	if *memProfile != "" {
		f, err := os.Create(*memProfile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		// get up-to-date statistics
		runtime.GC()

		if err := pprof.WriteHeapProfile(f); err != nil {
			panic(err)
		}
	}
}
