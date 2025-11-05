package main

import (
	"fmt"
	"log"
	"net/http"
)

// handler returns a simple greeting message
func handler(w http.ResponseWriter, r *http.Request) {
	message := "Hello from Go! 🚀\n"
	fmt.Fprintf(w, message)
	log.Printf("Served request from %s", r.RemoteAddr)
}

// healthHandler provides a health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/health", healthHandler)
	
	port := ":8080"
	log.Printf("Go service starting on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
