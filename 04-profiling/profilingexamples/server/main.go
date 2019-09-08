// Starts this program with
//
// go run main.go
//
// In other tab, collect the profile with
//
// go tool pprof http://localhost:9090/debug/pprof/profile

package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

// Fibonacci returns the n-th element of the Fibonacci sequence
func Fibonacci(n int) int {
	if n == 0 || n == 1 {
		return 1
	}

	return Fibonacci(n-2) + Fibonacci(n-1)
}

func main() {
	go func() {
		_ = Fibonacci(100)
	}()

	log.Fatal(http.ListenAndServe(":9090", nil))
}
