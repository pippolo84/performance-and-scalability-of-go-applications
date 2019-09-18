package main

import (
	"log"
	"net/http"
	"performance-and-scalability-of-go-applications/06-memory-and-gc/wordcounter/service"
)

func main() {
	// register the specific handler for the "/search" endpoint
	http.HandleFunc("/search", service.Search)

	// start the server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
