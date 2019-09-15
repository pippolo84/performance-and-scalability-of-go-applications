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
	"log"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
)

// Fibonacci returns the n-th element of the Fibonacci sequence
func Fibonacci(n int) int {
	if n == 0 || n == 1 {
		return 1
	}

	return Fibonacci(n-2) + Fibonacci(n-1)
}

func main() {
	cpuProfile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	fmt.Println(Fibonacci(100))
}
