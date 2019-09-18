package main

import (
	"log"
	"net/http"
	"performance-and-scalability-of-go-applications/06-memory-and-gc/wordcounter/solution/service"
	"runtime"

	_ "net/http/pprof"
)

func main() {
	runtime.MemProfileRate = 1

	// register the specific handler for the "/search" endpoint
	http.HandleFunc("/search", service.Search)

	// start the server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
